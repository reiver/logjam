import{u as o,F as m,T as c,B as d,M as f,t as u}from"./index-031006af.js";import{z as t,d as p,o as e}from"./index-acfd88c2.js";import"./metatagsLogo-ca8ad162.js";import"./helperAPI-632347b0.js";import"./logo-93ef2180.js";const x=t.object({room:t.string().min(1,"This field is required"),name:t.string().min(1,"This field is required")}),w=({params:{room:a}})=>{var l;const[s,i]=p(!1),r=o({defaultValues:{room:a,name:""},resolver:u(x)}),n=()=>{i(!0)};if(!s)return e("div",{class:"w-full flex justify-center items-center px-4",children:e("div",{class:"w-full flex justify-center items-center max-w-[500px] mx-auto mt-10 border rounded-md border-gray-300",children:e("form",{class:"flex flex-col w-full ",onSubmit:r.handleSubmit(n),children:[e("span",{className:"text-bold-12 text-black block text-center pt-5",children:"Join the meeting"}),e("hr",{className:"my-3"}),e("div",{className:"p-5 flex flex-col gap-3",children:[e("span",{class:"text-bold-12 text-gray-2",children:"Please enter your display name:"}),e(m,{className:"w-full",children:e(c,{label:"Display name",variant:"outlined",size:"small",...r.register("name"),error:!!r.formState.errors.name,helperText:(l=r.formState.errors.name)==null?void 0:l.message})}),e("div",{class:"flex gap-2 w-full",children:e(d,{type:"submit",variant:"contained",className:"w-full normal-case",sx:{textTransform:"none"},color:"primary",children:"Attend Live Show"})})]})]})})});if(s)return e(f,{params:{...r.getValues(),room:a}})};export{w as AudiencePage,w as default};
