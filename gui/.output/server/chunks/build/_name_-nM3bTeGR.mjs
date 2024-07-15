import { _ as _sfc_main$1 } from './PageHeader-maaT1_jD.mjs';
import { _ as __nuxt_component_2 } from './Badge-OGRmTTcV.mjs';
import { _ as __nuxt_component_3 } from './Button-CDYqYkPm.mjs';
import { a as useRoute } from './server.mjs';
import { ref, watch, unref, withCtx, openBlock, createBlock, createTextVNode, toDisplayString, createVNode, createCommentVNode, useSSRContext } from 'vue';
import { ssrRenderComponent, ssrInterpolate } from 'vue/server-renderer';
import 'tailwind-merge';
import './_plugin-vue_export-helper-1tPrXgE0.mjs';
import './nuxt-link-CgTv0-CK.mjs';
import '../runtime.mjs';
import 'node:http';
import 'node:https';
import 'events';
import 'https';
import 'http';
import 'net';
import 'tls';
import 'crypto';
import 'stream';
import 'url';
import 'zlib';
import 'buffer';
import 'fs';
import 'path';
import 'engine.io';
import 'socket.io';
import 'dockerode';
import 'node:fs';
import 'node:url';
import './Icon-Dm7V5HYH.mjs';
import '@iconify/vue/dist/offline';
import '@iconify/vue';
import './index-BETvu_7i.mjs';
import '../routes/renderer.mjs';
import 'vue-bundle-renderer/runtime';
import 'devalue';
import '@unhead/ssr';
import 'unhead';
import '@unhead/shared';
import 'vue-router';
import '@vueuse/core';

const _sfc_main = {
  __name: "[name]",
  __ssrInlineRender: true,
  setup(__props) {
    useRoute();
    const database = ref(null);
    watch(database, (newDatabase) => {
      database.value = newDatabase;
    });
    return (_ctx, _push, _parent, _attrs) => {
      const _component_UPageHeader = _sfc_main$1;
      const _component_UBadge = __nuxt_component_2;
      const _component_UButton = __nuxt_component_3;
      _push(`<!--[-->`);
      _push(ssrRenderComponent(_component_UPageHeader, {
        headline: unref(database) ? "" : "Loading...",
        class: "mb-4 sticky top-0 bg-white z-10 dark:bg-gray-900 pb-4"
      }, {
        title: withCtx((_, _push2, _parent2, _scopeId) => {
          if (_push2) {
            if (unref(database)) {
              _push2(`<h1 class="flex"${_scopeId}>${ssrInterpolate(unref(database).getName())} `);
              _push2(ssrRenderComponent(_component_UBadge, {
                label: unref(database).getType().label,
                color: unref(database).getType().color,
                class: "ml-2 h-6",
                variant: "soft"
              }, null, _parent2, _scopeId));
              _push2(`</h1>`);
            } else {
              _push2(`<!---->`);
            }
          } else {
            return [
              unref(database) ? (openBlock(), createBlock("h1", {
                key: 0,
                class: "flex"
              }, [
                createTextVNode(toDisplayString(unref(database).getName()) + " ", 1),
                createVNode(_component_UBadge, {
                  label: unref(database).getType().label,
                  color: unref(database).getType().color,
                  class: "ml-2 h-6",
                  variant: "soft"
                }, null, 8, ["label", "color"])
              ])) : createCommentVNode("", true)
            ];
          }
        }),
        description: withCtx((_, _push2, _parent2, _scopeId) => {
          if (_push2) {
            if (unref(database)) {
              _push2(`<div class="flex flex-wrap gap-4"${_scopeId}></div>`);
            } else {
              _push2(`<!---->`);
            }
          } else {
            return [
              unref(database) ? (openBlock(), createBlock("div", {
                key: 0,
                class: "flex flex-wrap gap-4"
              })) : createCommentVNode("", true)
            ];
          }
        }),
        _: 1
      }, _parent));
      if (unref(database)) {
        _push(`<div></div>`);
      } else {
        _push(`<div>`);
        _push(ssrRenderComponent(_component_UButton, {
          size: "xl",
          color: "white",
          square: "",
          block: "",
          variant: "ghost",
          loading: ""
        }, {
          default: withCtx((_, _push2, _parent2, _scopeId) => {
            if (_push2) {
              _push2(` Loading database... `);
            } else {
              return [
                createTextVNode(" Loading database... ")
              ];
            }
          }),
          _: 1
        }, _parent));
        _push(`</div>`);
      }
      _push(`<!--]-->`);
    };
  }
};
const _sfc_setup = _sfc_main.setup;
_sfc_main.setup = (props, ctx) => {
  const ssrContext = useSSRContext();
  (ssrContext.modules || (ssrContext.modules = /* @__PURE__ */ new Set())).add("pages/databases/[type]/[name].vue");
  return _sfc_setup ? _sfc_setup(props, ctx) : void 0;
};

export { _sfc_main as default };
//# sourceMappingURL=_name_-nM3bTeGR.mjs.map
