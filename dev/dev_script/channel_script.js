
const channel = require('./channel-operations.js')
const helper = require('./helper.js')
const config = require('./config.json')
const async = require('async')
var user_tokens  = {}
require('./config.js')

async function setup(params) {
    try {
        await  channel.joinChannel(channel_name,peers,username,orgname)
        await  channel.installChaincode(peers,chaincodeName,chaincodePath,chaincodeVersion,chaincodeType,username,orgname)
        await  channel.instantiateChaincode(peers,channel_name,chaincodeName,chaincodeVersion,fcn,chaincodeType,args,username,org_name,isUpgrade)
    } catch (err) {
        console.log(err)
    }
}

async function register(){
    console.log(config.orgs)
    for(var i in config.orgs){
        let response = await helper.getRegisteredUser(config.username, i,true);
        user_tokens[i] = response.secret
    }
}

async function createChannel(){
    let response = await channel.createChannel(config.channelName,config.channelConfigPath,config.username,Object.keys(config.orgs)[0])
    return response
}

// join channel
async function joinChannel(){
    var response
    for(var key in config.orgs){
        response = await  channel.joinChannel(config.channelName,config.orgs[key],config.username,key)
    }
    return response
}

//install chaincode
async function installChaincode(){
    var response
    for(var key in config.orgs){
        let response = await  channel.installChaincode(config.orgs[key],config.chaincodeName,config.chaincodePath,'v0','golang',config.username,key)
    }
    return response
}

async function instantiateChaincode( ){
    var allPeers=[]
    for(var key in config.orgs){
        allPeers.push(config.orgs[key][0])
    }
    let response = await  channel.instantiateChaincode(allPeers,config.channelName,config.chaincodeName,'v0','init','golang',config.instantiateArgs,config.username,Object.keys(config.orgs)[0],false)
    return response
}

async.waterfall(
    [
        function(cb){
            register().then(()=>{
                console.log('tokens',user_tokens)
                cb()
            })
        },
        function(cb){
	// setTimeout(async ()=>{
            createChannel().then((response)=>{
                console.log('create',response)
                cb()
            })
	// //},1000)
        },
        function(cb){
            setTimeout(async ()=>{
               joinChannel().then((response)=>{
                   console.log('join',response)
                   cb()
               })
            },1000)
        },
        function(cb){
	// setTimeout(async ()=>{
            installChaincode().then((response)=>{
                console.log('install',response)
                cb()
            })
	//},1000)
        },
        function(cb){
	 //setTimeout(async ()=>{
            instantiateChaincode().then((response)=>{
                console.log('instantiate',response)
                cb()
            })
	//},1000)
        }
    ],function(err){
        console.log('END...')
    })
