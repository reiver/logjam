var W=Object.defineProperty;var v=(o,t,a)=>t in o?W(o,t,{enumerable:!0,configurable:!0,writable:!0,value:a}):o[t]=a;var i=(o,t,a)=>(v(o,typeof t!="symbol"?t+"":t,a),a);import{a7 as w,a8 as f,c9 as A,aa as y,B as b,an as u,ca as g,f as C,_ as P,av as k,at as V,ap as _,ad as d,ae as p}from"./index-5efc63cb.js";import{C as x,a as E,G as T,b as D}from"./contract-appuri-9892f94f.browser.esm-13c4428f.js";import{C as I}from"./contract-interceptor-d7b164a7.browser.esm-7eabd2ea.js";let l=function(o){return o[o.Against=0]="Against",o[o.For=1]="For",o[o.Abstain=2]="Abstain",o}({});class R{constructor(t,a,r){i(this,"propose",d(async(t,a)=>{a||(a=[{toAddress:this.contractWrapper.address,nativeTokenValue:0,transactionData:"0x"}]);const r=a.map(s=>s.toAddress),n=a.map(s=>s.nativeTokenValue),c=a.map(s=>s.transactionData);return p.fromContractWrapper({contractWrapper:this.contractWrapper,method:"propose",args:[r,n,c,t],parse:s=>({id:this.contractWrapper.parseLogs("ProposalCreated",s==null?void 0:s.logs)[0].args.proposalId,receipt:s})})}));i(this,"vote",d((()=>{var t=this;return async function(a,r){let n=arguments.length>2&&arguments[2]!==void 0?arguments[2]:"";return await t.ensureExists(a),p.fromContractWrapper({contractWrapper:t.contractWrapper,method:"castVoteWithReason",args:[a,r,n]})}})()));i(this,"execute",d(async t=>{await this.ensureExists(t);const a=await this.get(t),r=a.executions.map(e=>e.toAddress),n=a.executions.map(e=>e.nativeTokenValue),c=a.executions.map(e=>e.transactionData),s=g(a.description);return p.fromContractWrapper({contractWrapper:this.contractWrapper,method:"execute",args:[r,n,c,s]})}));let n=arguments.length>3&&arguments[3]!==void 0?arguments[3]:{},c=arguments.length>4?arguments[4]:void 0,s=arguments.length>5?arguments[5]:void 0,e=arguments.length>6&&arguments[6]!==void 0?arguments[6]:new w(t,a,c,n,r);this._chainId=s,this.abi=f.parse(c||[]),this.contractWrapper=e,this.storage=r,this.metadata=new x(this.contractWrapper,A,this.storage),this.app=new E(this.contractWrapper,this.metadata,this.storage),this.encoder=new y(this.contractWrapper),this.estimator=new T(this.contractWrapper),this.events=new D(this.contractWrapper),this.interceptor=new I(this.contractWrapper)}get chainId(){return this._chainId}onNetworkUpdated(t){this.contractWrapper.updateSignerOrProvider(t)}getAddress(){return this.contractWrapper.address}async get(t){const r=(await this.getAll()).filter(n=>n.proposalId.eq(b.from(t)));if(r.length===0)throw new Error("proposal not found");return r[0]}async getAll(){const t=await this.contractWrapper.read("getAllProposals",[])??[];return(await Promise.all(t.map(r=>Promise.all([this.contractWrapper.read("state",[r.proposalId]),this.getProposalVotes(r.proposalId)])))).map((r,n)=>{let[c,s]=r;const e=t[n];return{proposalId:e.proposalId,proposer:e.proposer,description:e.description,startBlock:e.startBlock,endBlock:e.endBlock,state:c,votes:s,executions:e[3].map((m,h)=>({toAddress:e.targets[h],nativeTokenValue:m,transactionData:e.calldatas[h]}))}})}async getProposalVotes(t){const a=await this.contractWrapper.read("proposalVotes",[t]);return[{type:l.Against,label:"Against",count:a.againstVotes},{type:l.For,label:"For",count:a.forVotes},{type:l.Abstain,label:"Abstain",count:a.abstainVotes}]}async hasVoted(t,a){return a||(a=await this.contractWrapper.getSignerAddress()),this.contractWrapper.read("hasVoted",[t,await u(a)])}async canExecute(t){await this.ensureExists(t);const a=await this.get(t),r=a.executions.map(e=>e.toAddress),n=a.executions.map(e=>e.nativeTokenValue),c=a.executions.map(e=>e.transactionData),s=g(a.description);try{return await this.contractWrapper.callStatic().execute(r,n,c,s),!0}catch{return!1}}async balance(){const t=await this.contractWrapper.getProvider().getBalance(this.contractWrapper.address);return{name:"",symbol:"",decimals:18,value:t,displayValue:C(t,18)}}async balanceOfToken(t){const a=(await P(()=>import("./index-5efc63cb.js").then(n=>n.eG),["assets/index-5efc63cb.js","assets/index-36f85667.css"])).default,r=new k(await u(t),a,this.contractWrapper.getProvider());return await V(this.contractWrapper.getProvider(),t,await r.balanceOf(this.contractWrapper.address))}async ensureExists(t){try{await this.contractWrapper.read("state",[t])}catch{throw Error(`Proposal ${t} not found`)}}async settings(){const[t,a,r,n,c]=await Promise.all([this.contractWrapper.read("votingDelay",[]),this.contractWrapper.read("votingPeriod",[]),this.contractWrapper.read("token",[]),this.contractWrapper.read("quorumNumerator",[]),this.contractWrapper.read("proposalThreshold",[])]),s=await _(this.contractWrapper.getProvider(),r);return{votingDelay:t.toString(),votingPeriod:a.toString(),votingTokenAddress:r,votingTokenMetadata:s,votingQuorumFraction:n.toString(),proposalTokenThreshold:c.toString()}}async prepare(t,a,r){return p.fromContractWrapper({contractWrapper:this.contractWrapper,method:t,args:a,overrides:r})}async call(t,a,r){return this.contractWrapper.call(t,a,r)}}export{R as Vote};
