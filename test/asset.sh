

. common.sh


if [[ ! -f $CMD || ! -x $CMD ]]; then
	rm -f $CMD
	cp ../$CMD .
fi

if [[ ! -f $CONFIG ]]; then
	cp ../$CONFIG .
fi


$CMD wallet -c --name $WALLET --password $PASSWD
if (( $? != 0 )); then
	echo "wallet creation failed"
	exit 1
fi


