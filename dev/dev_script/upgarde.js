const channel = require('./channel-operations.js')
const helper = require('./helper.js')
const config = require('./config.json')
const async = require('async')
var user_tokens  = {}
require('./config.js')

var version = process.argv[2]


async function installChaincode(){
    var response
    for(var key in config.orgs){
        let response = await  channel.installChaincode(config.orgs[key],config.chaincodeName,config.chaincodePath,version,'golang',config.username,key)
    }
    return response
}


async function upgradeChaincode() {
    var allPeers=[]
    for(var key in config.orgs){
        allPeers.push(config.orgs[key][0])
    }
    let response = await  channel.instantiateChaincode(config.orgs['PackagerBosch'],config.channelName,config.chaincodeName,version,'init','golang',config.instantiateArgs,config.username,'PackagerBosch',true)
    return response
}


async.waterfall([
    function (cb){
        installChaincode().then((response)=>{
            console.log('install',response)
            cb()
        })
    },function (cb) {
        upgradeChaincode()
    }
])