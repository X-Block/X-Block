

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


output=$($CMD wallet -l --name $WALLET --password $PASSWD)
if (( $? != 0 )); then
	echo "wallet listing failed"
	exit 1
fi
programhash=$(echo "$output" | grep "program hash" | awk -F : '{print $2}')


output=$($CMD asset --reg --name XBlock --value 10000 --wallet $WALLET --password $PASSWD)
if (( $? != 0 )); then
	echo "asset registration failed"
	exit 1
fi
assetid=$(getHashFromOutput "$output")
echo "Asset ID: $assetid"

sleep 6


output=$($CMD asset --issue --asset $assetid --to $programhash --value 9999 --wallet $WALLET --password $PASSWD)
if (( $? != 0 )); then
	echo "asset issuance failed"
	exit 1
fi
issue=$(getHashFromOutput "$output")
echo "Issue: $issue"

sleep 6


output=$(./nodectl asset --transfer --asset $assetid --to $programhash --value 1 --wallet $WALLET --password $PASSWD)
if (( $? != 0 )); then
	echo "asset transfer failed"
	exit 1
fi
transfer=$(getHashFromOutput "$output")
echo "Transfer: $transfer"

echo PASS

exit 0
