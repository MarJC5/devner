import { a as __nuxt_component_3 } from "./Container-B3wCur7c.js";
import { _ as __nuxt_component_0 } from "./nuxt-link-CgTv0-CK.js";
import { n as nuxtLinkProps, u as useUI, c as getNuxtLinkProps, a as __nuxt_component_2 } from "./Button-CDYqYkPm.js";
import { defineComponent, useSlots, toRef, computed, mergeProps, unref, createSlots, withCtx, createVNode, openBlock, createBlock, createCommentVNode, renderSlot, createTextVNode, toDisplayString, useSSRContext } from "vue";
import { ssrRenderComponent, ssrRenderClass, ssrRenderSlot, ssrInterpolate } from "vue/server-renderer";
import "../server.mjs";
const getSlotChildrenText = (children) => children.map((node) => {
  if (!node.children || typeof node.children === "string") return node.children || "";
  else if (Array.isArray(node.children)) return getSlotChildrenText(node.children);
  else if (node.children.default) return getSlotChildrenText(node.children.default());
}).join("");
const _sfc_main = /* @__PURE__ */ defineComponent({
  ...{
    inheritAttrs: false
  },
  __name: "PageCard",
  __ssrInlineRender: true,
  props: {
    ...nuxtLinkProps,
    title: {
      type: String,
      default: void 0
    },
    description: {
      type: String,
      default: void 0
    },
    icon: {
      type: String,
      default: void 0
    },
    class: {
      type: [String, Object, Array],
      default: void 0
    },
    ui: {
      type: Object,
      default: () => ({})
    }
  },
  setup(__props) {
    const config = {
      wrapper: "relative group",
      to: "hover:ring-2 hover:ring-primary-500 dark:hover:ring-primary-400 hover:bg-gray-100/50 dark:hover:bg-gray-800/50",
      icon: {
        wrapper: "mb-6 flex",
        base: "w-10 h-10 flex-shrink-0 text-primary"
      },
      body: {
        base: "flex-1"
      },
      title: "text-gray-900 dark:text-white text-base font-semibold truncate flex items-center gap-1.5",
      description: "text-[15px] text-gray-500 dark:text-gray-400 mt-1"
    };
    const props = __props;
    const slots = useSlots();
    const { ui, attrs } = useUI("page.card", toRef(props, "ui"), config, toRef(props, "class"), true);
    const nuxtLinkBind = computed(() => getNuxtLinkProps(props));
    const ariaLabel = computed(() => (props.title || slots.title && getSlotChildrenText(slots.title()) || "Card link").trim());
    return (_ctx, _push, _parent, _attrs) => {
      const _component_UCard = __nuxt_component_3;
      const _component_NuxtLink = __nuxt_component_0;
      const _component_UIcon = __nuxt_component_2;
      _push(ssrRenderComponent(_component_UCard, mergeProps({
        class: [unref(ui).wrapper, _ctx.to && unref(ui).to]
      }, unref(attrs), { ui: unref(ui) }, _attrs), createSlots({
        default: withCtx((_, _push2, _parent2, _scopeId) => {
          if (_push2) {
            if (_ctx.to) {
              _push2(ssrRenderComponent(_component_NuxtLink, mergeProps({ "aria-label": unref(ariaLabel) }, unref(nuxtLinkBind), {
                class: "focus:outline-none",
                tabindex: "-1"
              }), {
                default: withCtx((_2, _push3, _parent3, _scopeId2) => {
                  if (_push3) {
                    _push3(`<span class="absolute inset-0" aria-hidden="true"${_scopeId2}></span>`);
                  } else {
                    return [
                      createVNode("span", {
                        class: "absolute inset-0",
                        "aria-hidden": "true"
                      })
                    ];
                  }
                }),
                _: 1
              }, _parent2, _scopeId));
            } else {
              _push2(`<!---->`);
            }
            if (__props.icon || _ctx.$slots.icon) {
              _push2(`<div class="${ssrRenderClass(unref(ui).icon.wrapper)}"${_scopeId}>`);
              ssrRenderSlot(_ctx.$slots, "icon", {}, () => {
                _push2(ssrRenderComponent(_component_UIcon, {
                  name: __props.icon,
                  class: unref(ui).icon.base
                }, null, _parent2, _scopeId));
              }, _push2, _parent2, _scopeId);
              _push2(`</div>`);
            } else {
              _push2(`<!---->`);
            }
            if (__props.title || _ctx.$slots.title) {
              _push2(`<p class="${ssrRenderClass(unref(ui).title)}"${_scopeId}>`);
              ssrRenderSlot(_ctx.$slots, "title", {}, () => {
                _push2(`${ssrInterpolate(__props.title)}`);
              }, _push2, _parent2, _scopeId);
              _push2(`</p>`);
            } else {
              _push2(`<!---->`);
            }
            if (__props.description || _ctx.$slots.description) {
              _push2(`<div class="${ssrRenderClass(unref(ui).description)}"${_scopeId}>`);
              ssrRenderSlot(_ctx.$slots, "description", {}, () => {
                _push2(`${ssrInterpolate(__props.description)}`);
              }, _push2, _parent2, _scopeId);
              _push2(`</div>`);
            } else {
              _push2(`<!---->`);
            }
            ssrRenderSlot(_ctx.$slots, "default", {}, null, _push2, _parent2, _scopeId);
          } else {
            return [
              _ctx.to ? (openBlock(), createBlock(_component_NuxtLink, mergeProps({
                key: 0,
                "aria-label": unref(ariaLabel)
              }, unref(nuxtLinkBind), {
                class: "focus:outline-none",
                tabindex: "-1"
              }), {
                default: withCtx(() => [
                  createVNode("span", {
                    class: "absolute inset-0",
                    "aria-hidden": "true"
                  })
                ]),
                _: 1
              }, 16, ["aria-label"])) : createCommentVNode("", true),
              __props.icon || _ctx.$slots.icon ? (openBlock(), createBlock("div", {
                key: 1,
                class: unref(ui).icon.wrapper
              }, [
                renderSlot(_ctx.$slots, "icon", {}, () => [
                  createVNode(_component_UIcon, {
                    name: __props.icon,
                    class: unref(ui).icon.base
                  }, null, 8, ["name", "class"])
                ])
              ], 2)) : createCommentVNode("", true),
              __props.title || _ctx.$slots.title ? (openBlock(), createBlock("p", {
                key: 2,
                class: unref(ui).title
              }, [
                renderSlot(_ctx.$slots, "title", {}, () => [
                  createTextVNode(toDisplayString(__props.title), 1)
                ])
              ], 2)) : createCommentVNode("", true),
              __props.description || _ctx.$slots.description ? (openBlock(), createBlock("div", {
                key: 3,
                class: unref(ui).description
              }, [
                renderSlot(_ctx.$slots, "description", {}, () => [
                  createTextVNode(toDisplayString(__props.description), 1)
                ])
              ], 2)) : createCommentVNode("", true),
              renderSlot(_ctx.$slots, "default")
            ];
          }
        }),
        _: 2
      }, [
        _ctx.$slots.header ? {
          name: "header",
          fn: withCtx((_, _push2, _parent2, _scopeId) => {
            if (_push2) {
              ssrRenderSlot(_ctx.$slots, "header", {}, null, _push2, _parent2, _scopeId);
            } else {
              return [
                renderSlot(_ctx.$slots, "header")
              ];
            }
          }),
          key: "0"
        } : void 0,
        _ctx.$slots.footer ? {
          name: "footer",
          fn: withCtx((_, _push2, _parent2, _scopeId) => {
            if (_push2) {
              ssrRenderSlot(_ctx.$slots, "footer", {}, null, _push2, _parent2, _scopeId);
            } else {
              return [
                renderSlot(_ctx.$slots, "footer")
              ];
            }
          }),
          key: "1"
        } : void 0
      ]), _parent));
    };
  }
});
const _sfc_setup = _sfc_main.setup;
_sfc_main.setup = (props, ctx) => {
  const ssrContext = useSSRContext();
  (ssrContext.modules || (ssrContext.modules = /* @__PURE__ */ new Set())).add("node_modules/@nuxt/ui-pro/components/page/PageCard.vue");
  return _sfc_setup ? _sfc_setup(props, ctx) : void 0;
};
export {
  _sfc_main as _
};
//# sourceMappingURL=PageCard-NhNi6yJq.js.map
