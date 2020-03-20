#!/bin/sh

function replacePrivateKey() {
  # sed on MacOSX does not support -i flag with a null extension. We will use
  # 't' for our back-up's extension and delete it at the end of the function
  # Copy the template to the file that will be modified to add the private key
  ARCH=$(uname -s | grep Darwin)
  if [ "$ARCH" == "Darwin" ]; then
    OPTS="-it"
  else
    OPTS="-i"
  fi
  cp docker-compose-e2e-template.yaml docker-compose-e2e.yaml

  # The next steps will replace the template's contents with the
  # actual values of the private key file names for the two CAs.
  CURRENT_DIR=$PWD
  cd crypto-config/peerOrganizations/packagerBosch.example.com/ca/
  PRIV_KEY=$(ls *_sk)
  cd "$CURRENT_DIR"
  sed $OPTS "s/CA0_PRIVATE_KEY/${PRIV_KEY}/g" docker-compose-e2e.yaml
  
  cd crypto-config/peerOrganizations/logistics.example.com/ca/
  PRIV_KEY=$(ls *_sk)
  cd "$CURRENT_DIR"
  sed $OPTS "s/CA1_PRIVATE_KEY/${PRIV_KEY}/g" docker-compose-e2e.yaml

  cd crypto-config/peerOrganizations/adc.example.com/ca/
  PRIV_KEY=$(ls *_sk)
  cd "$CURRENT_DIR"
  sed $OPTS "s/CA2_PRIVATE_KEY/${PRIV_KEY}/g" docker-compose-e2e.yaml

  cd crypto-config/peerOrganizations/ldc.example.com/ca/
  PRIV_KEY=$(ls *_sk)
  cd "$CURRENT_DIR"
  sed $OPTS "s/CA3_PRIVATE_KEY/${PRIV_KEY}/g" docker-compose-e2e.yaml

  cd crypto-config/peerOrganizations/dealer.example.com/ca/
  PRIV_KEY=$(ls *_sk)
  cd "$CURRENT_DIR"
  sed $OPTS "s/CA4_PRIVATE_KEY/${PRIV_KEY}/g" docker-compose-e2e.yaml

  cd crypto-config/peerOrganizations/retailer.example.com/ca/
  PRIV_KEY=$(ls *_sk)
  cd "$CURRENT_DIR"
  sed $OPTS "s/CA5_PRIVATE_KEY/${PRIV_KEY}/g" docker-compose-e2e.yaml

}

replacePrivateKey