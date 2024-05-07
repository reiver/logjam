var A=Object.defineProperty;var f=(i,c,r)=>c in i?A(i,c,{enumerable:!0,configurable:!0,writable:!0,value:r}):i[c]=r;var p=(i,c,r)=>(f(i,typeof c!="symbol"?c+"":c,r),r);import{az as W,a7 as C,a8 as T,aA as v,aa as y,ap as I,f as E,ad as m,an as k,ae as h,ao as R}from"./index-5efc63cb.js";import{h as $}from"./hasERC20Allowance-e9f79d1f.browser.esm-f4d54b7f.js";import{a as w}from"./marketplace-478946cc.browser.esm-155f1fb8.js";import{u as b}from"./QueryParams-af68a67b.browser.esm-0cfc440f.js";import{C as S,a as P,G as q,b as F}from"./contract-appuri-9892f94f.browser.esm-13c4428f.js";import{C as L,a as O}from"./contract-owner-dff73b10.browser.esm-5bbbb3cf.js";import{C as U}from"./contract-roles-93079777.browser.esm-8b287a0e.js";import{S as x}from"./erc-721-standard-28e4c2e3.browser.esm-000aff6f.js";import"./setErc20Allowance-96f7a033.browser.esm-6aaefd7c.js";import"./index-c1e6b325.js";import"./treeify-0bd2afd2.js";import"./assertEnabled-f866961a.browser.esm-186936d3.js";import"./erc-721-051c3774.browser.esm-d809507e.js";import"./drop-claim-conditions-6b858f0b.browser.esm-031e4d8f.js";const d=class d extends x{constructor(r,e,n){let o=arguments.length>3&&arguments[3]!==void 0?arguments[3]:{},t=arguments.length>4?arguments[4]:void 0,a=arguments.length>5?arguments[5]:void 0,s=arguments.length>6&&arguments[6]!==void 0?arguments[6]:new C(r,e,t,o,n);super(s,n,a);p(this,"wrap",m(async(r,e,n)=>{const[o,t,a]=await Promise.all([b(e,this.storage),this.toTokenStructList(r),k(n||await this.contractWrapper.getSignerAddress())]);return h.fromContractWrapper({contractWrapper:this.contractWrapper,method:"wrap",args:[t,o,a],parse:s=>{const u=this.contractWrapper.parseLogs("TokensWrapped",s==null?void 0:s.logs);if(u.length===0)throw new Error("TokensWrapped event not found");const l=u[0].args.tokenIdOfWrappedToken;return{id:l,receipt:s,data:()=>this.get(l)}}})}));p(this,"unwrap",m(async(r,e)=>{const n=await k(e||await this.contractWrapper.getSignerAddress());return h.fromContractWrapper({contractWrapper:this.contractWrapper,method:"unwrap",args:[r,n]})}));this.abi=T.parse(t||[]),this.metadata=new S(this.contractWrapper,v,this.storage),this.app=new P(this.contractWrapper,this.metadata,this.storage),this.roles=new U(this.contractWrapper,d.contractRoles),this.encoder=new y(this.contractWrapper),this.estimator=new q(this.contractWrapper),this.events=new F(this.contractWrapper),this.royalties=new L(this.contractWrapper,this.metadata),this.owner=new O(this.contractWrapper)}async getWrappedContents(r){const e=await this.contractWrapper.read("getWrappedContents",[r]),n=[],o=[],t=[];for(const a of e)switch(a.tokenType){case 0:{const s=await I(this.contractWrapper.getProvider(),a.assetContract);n.push({contractAddress:a.assetContract,quantity:E(a.totalAmount,s.decimals)});break}case 1:{o.push({contractAddress:a.assetContract,tokenId:a.tokenId});break}case 2:{t.push({contractAddress:a.assetContract,tokenId:a.tokenId,quantity:a.totalAmount.toString()});break}}return{erc20Tokens:n,erc721Tokens:o,erc1155Tokens:t}}async toTokenStructList(r){const e=[],n=this.contractWrapper.getProvider(),o=await this.contractWrapper.getSignerAddress();if(r.erc20Tokens)for(const t of r.erc20Tokens){const a=await R(n,t.quantity,t.contractAddress);if(!await $(this.contractWrapper,t.contractAddress,a))throw new Error(`ERC20 token with contract address "${t.contractAddress}" does not have enough allowance to transfer.

You can set allowance to the multiwrap contract to transfer these tokens by running:

await sdk.getToken("${t.contractAddress}").setAllowance("${this.getAddress()}", ${t.quantity});

`);e.push({assetContract:t.contractAddress,totalAmount:a,tokenId:0,tokenType:0})}if(r.erc721Tokens)for(const t of r.erc721Tokens){if(!await w(this.contractWrapper.getProvider(),this.getAddress(),t.contractAddress,t.tokenId,o))throw new Error(`ERC721 token "${t.tokenId}" with contract address "${t.contractAddress}" is not approved for transfer.

You can give approval the multiwrap contract to transfer this token by running:

await sdk.getNFTCollection("${t.contractAddress}").setApprovalForToken("${this.getAddress()}", ${t.tokenId});

`);e.push({assetContract:t.contractAddress,totalAmount:0,tokenId:t.tokenId,tokenType:1})}if(r.erc1155Tokens)for(const t of r.erc1155Tokens){if(!await w(this.contractWrapper.getProvider(),this.getAddress(),t.contractAddress,t.tokenId,o))throw new Error(`ERC1155 token "${t.tokenId}" with contract address "${t.contractAddress}" is not approved for transfer.

You can give approval the multiwrap contract to transfer this token by running:

await sdk.getEdition("${t.contractAddress}").setApprovalForAll("${this.getAddress()}", true);

`);e.push({assetContract:t.contractAddress,totalAmount:t.quantity,tokenId:t.tokenId,tokenType:2})}return e}async prepare(r,e,n){return h.fromContractWrapper({contractWrapper:this.contractWrapper,method:r,args:e,overrides:n})}async call(r,e,n){return this.contractWrapper.call(r,e,n)}};p(d,"contractRoles",W);let g=d;export{g as Multiwrap};
