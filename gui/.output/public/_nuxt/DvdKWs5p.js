import{H as y,I as g,y as f,J as u,m as v,L as k,M as h,o as r,f as $,w,c as d,O as i,Q as t,j as n,R as b,P as C,l as c}from"./D9CLLpjw.js";import{u as m}from"./DTL95yix.js";import{_ as j}from"./DlAUqK2U.js";const O={base:"",background:"bg-white dark:bg-gray-900",divide:"divide-y divide-gray-200 dark:divide-gray-800",ring:"ring-1 ring-gray-200 dark:ring-gray-800",rounded:"rounded-lg",shadow:"shadow",body:{base:"",background:"",padding:"px-4 py-5 sm:p-6"},header:{base:"",background:"",padding:"px-4 py-5 sm:px-6"},footer:{base:"",background:"",padding:"px-4 py-4 sm:px-6"}},S=y(g.ui.strategy,g.ui.card,O),A=f({inheritAttrs:!1,props:{as:{type:String,default:"div"},class:{type:[String,Object,Array],default:()=>""},ui:{type:Object,default:()=>({})}},setup(a){const{ui:e,attrs:s}=m("card",u(a,"ui"),S),o=v(()=>k(h(e.value.base,e.value.rounded,e.value.divide,e.value.ring,e.value.shadow,e.value.background),a.class));return{ui:e,attrs:s,cardClass:o}}});function B(a,e,s,o,l,p){return r(),$(C(a.$attrs.onSubmit?"form":a.as),b({class:a.cardClass},a.attrs),{default:w(()=>[a.$slots.header?(r(),d("div",{key:0,class:i([a.ui.header.base,a.ui.header.padding,a.ui.header.background])},[t(a.$slots,"header")],2)):n("",!0),a.$slots.default?(r(),d("div",{key:1,class:i([a.ui.body.base,a.ui.body.padding,a.ui.body.background])},[t(a.$slots,"default")],2)):n("",!0),a.$slots.footer?(r(),d("div",{key:2,class:i([a.ui.footer.base,a.ui.footer.padding,a.ui.footer.background])},[t(a.$slots,"footer")],2)):n("",!0)]),_:3},16,["class"])}const R=j(A,[["render",B]]),z=f({inheritAttrs:!1,__name:"PageGrid",props:{class:{type:[String,Object,Array],default:void 0},ui:{type:Object,default:()=>({})}},setup(a){const e={wrapper:"grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 gap-8"},s=a,{ui:o,attrs:l}=m("page.grid",u(s,"ui"),e,u(s,"class"),!0);return(p,P)=>(r(),d("div",b({class:c(o).wrapper},c(l)),[t(p.$slots,"default")],16))}});export{z as _,R as a};
