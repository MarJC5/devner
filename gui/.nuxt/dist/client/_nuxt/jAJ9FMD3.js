import{_ as N}from"./CFu4p6eu.js";import{n as B,u as L,c as _,a as j}from"./DTL95yix.js";import{a as A}from"./DvdKWs5p.js";import{y as P,W as V,J as f,m as y,o as i,f as g,X as O,w as l,Q as o,R as h,l as t,j as c,c as u,O as d,b as U,d as k,t as b,a as I}from"./D9CLLpjw.js";const p=a=>a.map(s=>{if(!s.children||typeof s.children=="string")return s.children||"";if(Array.isArray(s.children))return p(s.children);if(s.children.default)return p(s.children.default())}).join(""),R=I("span",{class:"absolute inset-0","aria-hidden":"true"},null,-1),Q=P({inheritAttrs:!1,__name:"PageCard",props:{...B,title:{type:String,default:void 0},description:{type:String,default:void 0},icon:{type:String,default:void 0},class:{type:[String,Object,Array],default:void 0},ui:{type:Object,default:()=>({})}},setup(a){const s={wrapper:"relative group",to:"hover:ring-2 hover:ring-primary-500 dark:hover:ring-primary-400 hover:bg-gray-100/50 dark:hover:bg-gray-800/50",icon:{wrapper:"mb-6 flex",base:"w-10 h-10 flex-shrink-0 text-primary"},body:{base:"flex-1"},title:"text-gray-900 dark:text-white text-base font-semibold truncate flex items-center gap-1.5",description:"text-[15px] text-gray-500 dark:text-gray-400 mt-1"},n=a,m=V(),{ui:r,attrs:v}=L("page.card",f(n,"ui"),s,f(n,"class"),!0),x=y(()=>_(n)),$=y(()=>(n.title||m.title&&p(m.title())||"Card link").trim());return(e,T)=>{const S=N,w=j,C=A;return i(),g(C,h({class:[t(r).wrapper,e.to&&t(r).to]},t(v),{ui:t(r)}),O({default:l(()=>[e.to?(i(),g(S,h({key:0,"aria-label":t($)},t(x),{class:"focus:outline-none",tabindex:"-1"}),{default:l(()=>[R]),_:1},16,["aria-label"])):c("",!0),a.icon||e.$slots.icon?(i(),u("div",{key:1,class:d(t(r).icon.wrapper)},[o(e.$slots,"icon",{},()=>[U(w,{name:a.icon,class:d(t(r).icon.base)},null,8,["name","class"])])],2)):c("",!0),a.title||e.$slots.title?(i(),u("p",{key:2,class:d(t(r).title)},[o(e.$slots,"title",{},()=>[k(b(a.title),1)])],2)):c("",!0),a.description||e.$slots.description?(i(),u("div",{key:3,class:d(t(r).description)},[o(e.$slots,"description",{},()=>[k(b(a.description),1)])],2)):c("",!0),o(e.$slots,"default")]),_:2},[e.$slots.header?{name:"header",fn:l(()=>[o(e.$slots,"header")]),key:"0"}:void 0,e.$slots.footer?{name:"footer",fn:l(()=>[o(e.$slots,"footer")]),key:"1"}:void 0]),1040,["class","ui"])}}});export{Q as _};