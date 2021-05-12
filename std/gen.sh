#!/bin/bash

std_dir="$(dirname "$0")"

infer_package_name() {
	# Convert the filename $1 to a package name. We munge the name as follows:
	#
	# 1. strip off the capnp file extension and dirname
	# 2. remove dashes
	# 3. convert '+' to 'x'. This is really just for c++.capnp, but it's not
	#    any easier to special case it.
	printf '%s' "$(basename $1)" | sed 's/\.capnp$// ; s/-//g ; s/+/x/g'
}

gen_annotated_schema() {
	# Copy the schema from file "$1" to the std/capnp directory, and add
	# appropriate $Go annotations.
	infile="$1"
	outfile="$std_dir/capnp/$(basename "$infile")"
	package_name="$(infer_package_name "$outfile")"
	cat "$infile" - > "$outfile" << EOF
using Go = import "/go.capnp";
\$Go.package("$package_name");
\$Go.import("capnproto.org/go/capnp/v3/std/capnp/$package_name");
EOF
}

gen_go_src() {
	# Generate go source code from the schema file $1. Create the package
	# directory if necessary.
	file="$1"
        filedir="$(dirname "$file")"
	package_name="$(infer_package_name "$file")"
	mkdir -p "$filedir/$package_name" && \
        	capnp compile --no-standard-import -I"$std_dir" -ogo:"$filedir/$package_name" --src-prefix="$filedir" "$file"
}

usage() {
	echo "Usage:"
	echo ""
	echo "    $0 import <path/to/capnp/c++/src/capnp>"
	echo "    $0 compile    # Generate go source files"
	echo "    $0 clean-go   # Remove go source files"
	echo "    $0 clean-all  # Remove go source files and imported schemas"
}

# do_* implements the corresponding subcommand described in usage's output.
do_import() {
	input_dir="$1"
	for file in "$input_dir"/*.capnp; do
		gen_annotated_schema "$file" || return 1
	done
}

do_compile() {
	for file in "$std_dir"/*.capnp "$std_dir"/capnp/*.capnp; do
		gen_go_src "$file" || return 1
	done
}

do_clean_go() {
	find "$std_dir" -name '*.capnp.go' -delete
	find "$std_dir" -type d -empty -delete
}

do_clean_all() {
	do_clean_go
	find "$std_dir/capnp" -name '*.capnp' -delete
}

eq_or_usage() {
	# If "$1" is not equal to "$2", call usage and exit.
	if [ ! $1 = $2 ] ; then
		usage
		exit 1
	fi
}

case "$1" in
	import)    eq_or_usage $# 2; do_import "$2" ;;
	compile)   eq_or_usage $# 1; do_compile ;;
	clean-go)  eq_or_usage $# 1; do_clean_go ;;
	clean-all) eq_or_usage $# 1; do_clean_all ;;
	*) usage; exit 1 ;;
esac
