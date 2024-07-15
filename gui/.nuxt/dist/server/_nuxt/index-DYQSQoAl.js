import { _ as _sfc_main$1 } from "./PageHeader-maaT1_jD.js";
import { mergeProps, useSSRContext } from "vue";
import { ssrRenderComponent } from "vue/server-renderer";
import { _ as _export_sfc } from "./_plugin-vue_export-helper-1tPrXgE0.js";
import "./Button-CDYqYkPm.js";
import "./nuxt-link-CgTv0-CK.js";
import "ufo";
import "../server.mjs";
import "ofetch";
import "#internal/nuxt/paths";
import "hookable";
import "unctx";
import "h3";
import "unhead";
import "@unhead/shared";
import "vue-router";
import "radix3";
import "defu";
import "devalue";
import "@vueuse/core";
import "klona";
import "tailwind-merge";
import "ohash";
import "./Icon-Dm7V5HYH.js";
import "@iconify/vue/dist/offline";
import "@iconify/vue";
import "./index-BETvu_7i.js";
const _sfc_main = {};
function _sfc_ssrRender(_ctx, _push, _parent, _attrs) {
  const _component_UPageHeader = _sfc_main$1;
  _push(ssrRenderComponent(_component_UPageHeader, mergeProps({
    title: "Projects",
    class: "sticky top-0 bg-white z-10 dark:bg-gray-800 mb-4 pb-4"
  }, _attrs), null, _parent));
}
const _sfc_setup = _sfc_main.setup;
_sfc_main.setup = (props, ctx) => {
  const ssrContext = useSSRContext();
  (ssrContext.modules || (ssrContext.modules = /* @__PURE__ */ new Set())).add("pages/projects/index.vue");
  return _sfc_setup ? _sfc_setup(props, ctx) : void 0;
};
const index = /* @__PURE__ */ _export_sfc(_sfc_main, [["ssrRender", _sfc_ssrRender]]);
export {
  index as default
};
//# sourceMappingURL=index-DYQSQoAl.js.map
