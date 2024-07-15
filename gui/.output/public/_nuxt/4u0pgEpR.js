import{H as f,I as u,y as h,J as m,m as y,L as b,M as w,o as x,c as k,Q as S,d as v,t as $,R as A}from"./D9CLLpjw.js";import{u as C,b as O}from"./DTL95yix.js";import{_ as j}from"./DlAUqK2U.js";const N={base:"inline-flex items-center",rounded:"rounded-md",font:"font-medium",size:{xs:"text-xs px-1.5 py-0.5",sm:"text-xs px-2 py-1",md:"text-sm px-2 py-1",lg:"text-sm px-2.5 py-1.5"},color:{white:{solid:"ring-1 ring-inset ring-gray-300 dark:ring-gray-700 text-gray-900 dark:text-white bg-white dark:bg-gray-900"},gray:{solid:"ring-1 ring-inset ring-gray-300 dark:ring-gray-700 text-gray-700 dark:text-gray-200 bg-gray-50 dark:bg-gray-800"},black:{solid:"text-white dark:text-gray-900 bg-gray-900 dark:bg-white"}},variant:{solid:"bg-{color}-500 dark:bg-{color}-400 text-white dark:text-gray-900",outline:"text-{color}-500 dark:text-{color}-400 ring-1 ring-inset ring-{color}-500 dark:ring-{color}-400",soft:"bg-{color}-50 dark:bg-{color}-400 dark:bg-opacity-10 text-{color}-500 dark:text-{color}-400",subtle:"bg-{color}-50 dark:bg-{color}-400 dark:bg-opacity-10 text-{color}-500 dark:text-{color}-400 ring-1 ring-inset ring-{color}-500 dark:ring-{color}-400 ring-opacity-25 dark:ring-opacity-25"},default:{size:"sm",variant:"solid",color:"primary"}},c=f(u.ui.strategy,u.ui.badge,N),E=h({inheritAttrs:!1,props:{size:{type:String,default:()=>c.default.size,validator(s){return Object.keys(c.size).includes(s)}},color:{type:String,default:()=>c.default.color,validator(s){return[...u.ui.colors,...Object.keys(c.color)].includes(s)}},variant:{type:String,default:()=>c.default.variant,validator(s){return[...Object.keys(c.variant),...Object.values(c.color).flatMap(t=>Object.keys(t))].includes(s)}},label:{type:[String,Number],default:null},class:{type:[String,Object,Array],default:()=>""},ui:{type:Object,default:()=>({})}},setup(s){const{ui:t,attrs:e}=C("badge",m(s,"ui"),c),{size:o,rounded:r}=O({ui:t,props:s}),a=y(()=>{var n,l;const i=((l=(n=t.value.color)==null?void 0:n[s.color])==null?void 0:l[s.variant])||t.value.variant[s.variant];return b(w(t.value.base,t.value.font,r.value,t.value.size[o.value],i==null?void 0:i.replaceAll("{color}",s.color)),s.class)});return{attrs:e,badgeClass:a}}});function P(s,t,e,o,r,a){return x(),k("span",A({class:s.badgeClass},s.attrs),[S(s.$slots,"default",{},()=>[v($(s.label),1)])],16)}const L=j(E,[["render",P]]);function T({onlyFirst:s=!1}={}){const t=["[\\u001B\\u009B][[\\]()#;?]*(?:(?:(?:(?:;[-a-zA-Z\\d\\/#&.:=?%@~_]+)*|[a-zA-Z\\d]+(?:;[-a-zA-Z\\d\\/#&.:=?%@~_]*)*)?\\u0007)","(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PR-TZcf-ntqry=><~]))"].join("|");return new RegExp(t,s?void 0:"g")}const z=T();function g(s){if(typeof s!="string")throw new TypeError(`Expected a \`string\`, got \`${typeof s}\``);return s.replace(z,"")}class d{constructor(t){this.id=t.Id||0,this.created=t.Created||"",this.name=t.Name||"",this.status=t.State.Status||"",this.pid=t.State.Pid||0,this.startedAt=t.State.StartedAt||"",this.image=t.Image||"",this.ports=t.NetworkSettings.Ports||{},this.labels=t.Config.Labels||{},this.network=t.NetworkSettings.Networks||{},this.project=this.labels["com.docker.compose.project"]||"",this.loading={start:!1,restart:!1,stop:!1,remove:!1,rebuild:!1}}static async logs(t,e={}){const o=await $fetch(`/api/containers/${t}/logs`,{method:"POST",body:JSON.stringify(e),headers:{"Content-Type":"application/json"}});o.on("data",r=>r.toString("utf8")),o.on("end",()=>{console.log("Stream ended")}),o.on("error",r=>{console.error("Stream error:",r)})}static async fetchContainer(t){const e=await $fetch(`/api/containers/${t}/details`,{method:"GET"});return new d(e)}static async fetchContainerByName(t){try{return(await d.all()).find(o=>o.name.includes(`/${t}`))}catch(e){throw console.error("Error fetching container by name:",e),e}}static async all(){const t=await $fetch("/api/containers",{method:"GET"});return await Promise.all(t.map(async o=>{const r=await d.fetchContainer(o.Id);if(r.getNetworks().devner&&!r.getName().toLowerCase().includes("gui"))return r})).then(o=>o.filter(r=>r!==void 0))}static async allOthers(){const t=await $fetch("/api/containers",{method:"GET"}),o=(await Promise.all(t.map(async a=>{const i=await d.fetchContainer(a.Id);if(!i.getNetworks().devner)return i})).then(a=>a.filter(i=>i!==void 0))).reduce((a,i)=>{const n=i.project;return a[n]||(a[n]=[]),a[n].push(i),a},{});return Object.entries(o).map(([a,i])=>({label:a.charAt(0).toUpperCase()+a.slice(1),icon:"i-heroicons-cube-transparent",containers:i}))}formatDate(t){return new Date(t).toLocaleString()}async performAction(t){this.loading[t]=!0,console.log(`Performing action ${t} on container with ID: ${this.id}`),await $fetch(`/api/containers/${this.id}/${t}`,{method:"POST"}),this.loading[t]=!1,await this.updateStatus()}async start(){await this.performAction("start")}async restart(){await this.performAction("restart")}async stop(){await this.performAction("stop")}async remove(){await this.performAction("remove")}async rebuild(){await this.performAction("rebuild")}async updateStatus(){const t=await d.fetchContainer(this.id);this.status=t.status,this.pid=t.pid,this.startedAt=t.startedAt}async cmd(t,e=!1,o=!0){try{const r=await $fetch(`/api/containers/${this.id}/exec`,{method:"POST",body:JSON.stringify({command:t}),headers:{"Content-Type":"application/json"}});if(r.status==="error")throw new Error(r.message);const a=r.stdout.split(`
`);return e?this.parseTableOutput(a):o?a.map(i=>g(i).replace(/[^\x20-\x7E]/g,"")).filter(i=>i):a}catch(r){throw console.error("Error executing command in container:",r),r}}parseTableOutput(t){const[e,...o]=t,r=e.split("	").map(a=>a.trim());return o.map(a=>{const i=a.split("	").map(l=>l.trim());let n={};return r.forEach((l,p)=>{n[g(l).replace(/[^\x20-\x7E]/g,"")]=i[p]||""}),n})}async getDirectories(t="/"){const e=await this.cmd(`ls -a ${t}`);return{filesAndDirs:e,path:t,count:e.length}}getId(){return this.id}getCreated(){return this.formatDate(this.created)}getName(){let t=this.name.replace("/","").replace(/[-_]/g," ").replace(/\b\w/g,e=>e.toUpperCase());return t=t.replace(/ Devner$/,""),t}getStatus(t="capitalize"){return t==="capitalize"?this.status.charAt(0).toUpperCase()+this.status.slice(1):this.status}getImage(){return this.image}getPorts(t=!0){const e=[];for(const[o,r]of Object.entries(this.ports))r&&r.forEach(a=>{t?e.push(`${a.HostPort}/${o.split("/")[1]}`):e.push(a.HostPort)});return e}getLabels(){return this.labels}getNetworks(){return this.network}isRunning(){return this.status==="running"}isLoading(t){return this.loading[t]}}export{d as C,L as _};
