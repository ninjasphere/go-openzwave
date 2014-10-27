die()
{
    echo "$*" 1>&2
    exit 1
}

scenario()
{
    record()
    {
	local name="$1"
	local dir="scenarios/$name"
	test -d "$dir" || mkdir -p "$dir"
	
	( 
	    cd $dir &&
	    ../../go/bin/zwcli -configDir ../../../go-openzwave/openzwave/config -save -monitor -debug -logFileName=log 2>&1 | tee output
	) || exit $?
    }

    "$@"
}

scenario "$@"