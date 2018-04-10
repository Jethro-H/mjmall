#!/bin/sh

CHANNEL_NAME=$1
RDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/mjmall.com/orderers/orderer.mjmall.com/msp/tlscacerts/tlsca.mjmall.com-cert.pem

. ./utils.sh

createChannel() {
#        setGlobals 0 1

        if [ -z "$CORE_PEER_TLS_ENABLED" -o "$CORE_PEER_TLS_ENABLED" = "false" ]; then
		echo "ttt"
                set -x
                peer channel create -o orderer.mjmall.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/channel.tx >&log.txt
#                res=$?
                set +x
        else
		echo $CHANNEL_NAME
                set -x
                peer channel create -o orderer.mjmall.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/channel.tx --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA >&log.txt
#                res=$?
                set +x
        fi
#        cat log.txt
#        verifyResult $res "Channel creation failed"
        echo "===================== Channel \"$CHANNEL_NAME\" is created successfully ===================== "
        echo
}

#joinChannel () {
#        for org in 1 2; do
#            for peer in 0 1; do
#                joinChannelWithRetry $peer $org
#                echo "===================== peer${peer}.org${org} joined on the channel \"$CHANNEL_NAME\" ===================== "
#                sleep $DELAY
#                echo
#            done
#        done
#}
#
#updateAnchorPeers() {
#  PEER=$1
#  ORG=$2
#  setGlobals $PEER $ORG
#
#  if [ -z "$CORE_PEER_TLS_ENABLED" -o "$CORE_PEER_TLS_ENABLED" = "false" ]; then
#                set -x
#                peer channel update -o orderer.mjmall.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/${CORE_PEER_LOCALMSPID}anchors.tx >&log.txt
#                res=$?
#                set +x
#  else
#                set -x
#                peer channel update -o orderer.mjmall.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/${CORE_PEER_LOCALMSPID}anchors.tx --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA >&log.txt
#                res=$?
#                set +x
#  fi
#        cat log.txt
#        verifyResult $res "Anchor peer update failed"
#        echo "===================== Anchor peers for org \"$CORE_PEER_LOCALMSPID\" on \"$CHANNEL_NAME\" is updated successfully ===================== "
#        sleep $DELAY
#        echo
#}

createChannel

#joinChannel

#updateAnchorPeers 0 1

#updateAnchorPeers 0 2
