var o=Object.defineProperty;var i=(e,r,t)=>r in e?o(e,r,{enumerable:!0,configurable:!0,writable:!0,value:t}):e[r]=t;var a=(e,r,t)=>(i(e,typeof r!="symbol"?r+"":r,t),t);import{aj as n,a8 as c,a9 as p}from"./index-f6fefe02.js";import{a as h}from"./erc-721-051c3774.browser.esm-6687f9f0.js";class d{constructor(r,t,s){a(this,"transfer",c(async(r,t)=>this.erc721.transfer.prepare(r,t)));a(this,"setApprovalForAll",c(async(r,t)=>this.erc721.setApprovalForAll.prepare(r,t)));a(this,"setApprovalForToken",c(async(r,t)=>p.fromContractWrapper({contractWrapper:this.contractWrapper,method:"approve",args:[await n(r),t]})));this.contractWrapper=r,this.storage=t,this.erc721=new h(this.contractWrapper,this.storage,s),this._chainId=s}get chainId(){return this._chainId}onNetworkUpdated(r){this.contractWrapper.updateSignerOrProvider(r)}getAddress(){return this.contractWrapper.address}async getAll(r){return this.erc721.getAll(r)}async getOwned(r,t){return r&&(r=await n(r)),this.erc721.getOwned(r,t)}async getOwnedTokenIds(r){return r&&(r=await n(r)),this.erc721.getOwnedTokenIds(r)}async totalSupply(){return this.erc721.totalCirculatingSupply()}async get(r){return this.erc721.get(r)}async ownerOf(r){return this.erc721.ownerOf(r)}async balanceOf(r){return this.erc721.balanceOf(r)}async balance(){return this.erc721.balance()}async isApproved(r,t){return this.erc721.isApproved(r,t)}}export{d as S};
