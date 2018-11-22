

. common.sh


if [[ ! -f $CMD || ! -x $CMD ]]; then
	rm -f $CMD
	cp ../$CMD .
fi

