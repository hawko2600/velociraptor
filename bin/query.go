/*
   Velociraptor - Hunting Evil
   Copyright (C) 2019 Velocidex Innovations.

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/Velocidex/ordereddict"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	artifacts "www.velocidex.com/golang/velociraptor/artifacts"
	"www.velocidex.com/golang/velociraptor/file_store/csv"
	"www.velocidex.com/golang/velociraptor/reporting"
	"www.velocidex.com/golang/velociraptor/uploads"
	vql_subsystem "www.velocidex.com/golang/velociraptor/vql"
	"www.velocidex.com/golang/vfilter"
)

var (
	// Command line interface for VQL commands.
	query   = app.Command("query", "Run a VQL query")
	queries = query.Arg("queries", "The VQL Query to run.").
		Required().Strings()

	rate = app.Flag("ops_per_second", "Rate of execution").
		Default("1000000").Float64()
	format = query.Flag("format", "Output format to use (text,json,csv,jsonl).").
		Default("json").Enum("text", "json", "csv", "jsonl")
	dump_dir = query.Flag("dump_dir", "Directory to dump output files.").
			Default(".").String()

	env_map = app.Flag("env", "Environment for the query.").
		StringMap()

	max_wait = app.Flag("max_wait", "Maximum time to queue results.").
			Default("10").Int()

	explain        = app.Command("explain", "Explain the output from a plugin")
	explain_plugin = explain.Arg("plugin", "Plugin to explain").Required().String()
)

func outputJSON(ctx context.Context,
	scope *vfilter.Scope,
	vql *vfilter.VQL,
	out io.Writer) {
	result_chan := vfilter.GetResponseChannel(vql, ctx, scope, 10, *max_wait)
	for {
		result, ok := <-result_chan
		if !ok {
			return
		}
		out.Write(result.Payload)
	}
}

func outputJSONL(ctx context.Context,
	scope *vfilter.Scope,
	vql *vfilter.VQL,
	out io.Writer) {
	result_chan := vfilter.GetResponseChannel(vql, ctx, scope, 10, *max_wait)
rows:
	for {
		result, ok := <-result_chan
		if !ok {
			return
		}

		result_array := []json.RawMessage{}
		err := json.Unmarshal(result.Payload, &result_array)
		if err != nil {
			continue rows
		}

		for _, item := range result_array {
			// Decode the row into an ordered dict to maintain ordering.
			row := ordereddict.NewDict()
			err = json.Unmarshal(item, row)
			if err != nil {
				continue
			}

			// Re-serialize it as compact json.
			serialized, err := json.Marshal(row)
			if err != nil {
				continue
			}

			out.Write(serialized)

			// Separate lines with \n
			out.Write([]byte("\n"))
		}
	}
}

func outputCSV(ctx context.Context,
	scope *vfilter.Scope,
	vql *vfilter.VQL,
	out io.Writer) {
	result_chan := vfilter.GetResponseChannel(vql, ctx, scope, 10, *max_wait)

	csv_writer := csv.GetCSVAppender(
		scope, &StdoutWrapper{out}, true /* write_headers */)
	defer csv_writer.Close()

	for result := range result_chan {
		payload := []map[string]interface{}{}
		err := json.Unmarshal(result.Payload, &payload)
		kingpin.FatalIfError(err, "outputCSV")

		for _, row := range payload {
			row_dict := ordereddict.NewDict()
			for _, column := range result.Columns {
				value, pres := row[column]
				if pres {
					row_dict.Set(column, value)
				}
			}

			csv_writer.Write(row_dict)
		}
	}

}

func doQuery() {
	config_obj := get_config_or_default()
	repository, err := artifacts.GetGlobalRepository(config_obj)
	kingpin.FatalIfError(err, "Artifact GetGlobalRepository ")

	if *artifact_definitions_dir != "" {
		repository.LoadDirectory(*artifact_definitions_dir)
	}

	var acl_manager vql_subsystem.ACLManager = vql_subsystem.NullACLManager{}
	if *run_as != "" {
		acl_manager = vql_subsystem.NewServerACLManager(config_obj, *run_as)
	}

	env := ordereddict.NewDict().
		Set("config", config_obj.Client).
		Set("server_config", config_obj).
		Set("$uploader", &uploads.FileBasedUploader{
			UploadDir: *dump_dir,
		}).

		// Running on the commandline has no ACL restrictions.
		Set(vql_subsystem.ACL_MANAGER_VAR, acl_manager).
		Set(vql_subsystem.CACHE_VAR, vql_subsystem.NewScopeCache())

	if env_map != nil {
		for k, v := range *env_map {
			env.Set(k, v)
		}
	}

	scope := artifacts.MakeScope(repository).AppendVars(env)
	defer scope.Close()

	// Install throttler into the scope.
	vfilter.InstallThrottler(scope, vfilter.NewTimeThrottler(float64(*rate)))

	ctx := InstallSignalHandler(scope)

	AddLogger(scope, get_config_or_default())
	if *trace_vql_flag {
		scope.Tracer = log.New(os.Stderr, "VQL Trace: ", log.Lshortfile)
	}
	for _, query := range *queries {
		vql, err := vfilter.Parse(query)
		if err != nil {
			kingpin.FatalIfError(err, "Unable to parse VQL Query")
		}

		switch *format {
		case "text":
			table := reporting.EvalQueryToTable(ctx, scope, vql, os.Stdout)
			table.Render()
		case "json":
			outputJSON(ctx, scope, vql, os.Stdout)

		case "jsonl":
			outputJSONL(ctx, scope, vql, os.Stdout)

		case "csv":
			outputCSV(ctx, scope, vql, os.Stdout)
		}
	}
}

func doExplain(plugin string) {
	result := ordereddict.NewDict()
	type_map := vfilter.NewTypeMap()
	scope := vql_subsystem.MakeScope()
	defer scope.Close()

	pslist_info, pres := scope.Info(type_map, plugin)
	if pres {
		result.Set(plugin+"_info", pslist_info)
		result.Set("type_map", type_map)
	}

	s, err := json.MarshalIndent(result, "", " ")
	if err == nil {
		os.Stdout.Write(s)
	}
}

func init() {
	command_handlers = append(command_handlers, func(command string) bool {
		switch command {
		case explain.FullCommand():
			doExplain(*explain_plugin)

		case query.FullCommand():
			doQuery()

		default:
			return false
		}
		return true
	})
}
