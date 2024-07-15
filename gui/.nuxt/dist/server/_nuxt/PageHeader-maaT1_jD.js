import { u as useUI, a as __nuxt_component_2, _ as __nuxt_component_3 } from "./Button-CDYqYkPm.js";
import { defineComponent, toRef, mergeProps, unref, useSSRContext } from "vue";
import { ssrRenderAttrs, ssrRenderClass, ssrRenderSlot, ssrInterpolate, ssrRenderComponent, ssrRenderList } from "vue/server-renderer";
const _sfc_main = /* @__PURE__ */ defineComponent({
  ...{
    inheritAttrs: false
  },
  __name: "PageHeader",
  __ssrInlineRender: true,
  props: {
    headline: {
      type: String,
      default: void 0
    },
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
    links: {
      type: Array,
      default: () => []
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
      wrapper: "relative border-b border-gray-200 dark:border-gray-800 py-8",
      container: "flex flex-col lg:flex-row lg:items-center lg:justify-between",
      headline: "mb-3 text-sm/6 font-semibold text-primary flex items-center gap-1.5",
      title: "text-3xl sm:text-4xl font-bold text-gray-900 dark:text-white tracking-tight",
      description: "mt-4 text-lg text-gray-500 dark:text-gray-400",
      icon: {
        wrapper: "flex",
        base: "w-10 h-10 flex-shrink-0 text-primary"
      },
      links: "flex flex-wrap items-center gap-1.5 mt-4 lg:mt-0"
    };
    const props = __props;
    const { ui, attrs } = useUI("page.header", toRef(props, "ui"), config, toRef(props, "class"), true);
    return (_ctx, _push, _parent, _attrs) => {
      var _a;
      const _component_UIcon = __nuxt_component_2;
      const _component_UButton = __nuxt_component_3;
      _push(`<div${ssrRenderAttrs(mergeProps({
        class: unref(ui).wrapper
      }, unref(attrs), _attrs))}>`);
      if (__props.headline || _ctx.$slots.headline) {
        _push(`<div class="${ssrRenderClass(unref(ui).headline)}">`);
        ssrRenderSlot(_ctx.$slots, "headline", {}, () => {
          _push(`${ssrInterpolate(__props.headline)}`);
        }, _push, _parent);
        _push(`</div>`);
      } else {
        _push(`<!---->`);
      }
      _push(`<div class="flex flex-col lg:flex-row items-start gap-6">`);
      if (__props.icon || _ctx.$slots.icon) {
        _push(`<div class="${ssrRenderClass(unref(ui).icon.wrapper)}">`);
        ssrRenderSlot(_ctx.$slots, "icon", {}, () => {
          _push(ssrRenderComponent(_component_UIcon, {
            name: __props.icon,
            class: unref(ui).icon.base
          }, null, _parent));
        }, _push, _parent);
        _push(`</div>`);
      } else {
        _push(`<!---->`);
      }
      _push(`<div class="flex-1"><div class="${ssrRenderClass(unref(ui).container)}"><h1 class="${ssrRenderClass(unref(ui).title)}">`);
      ssrRenderSlot(_ctx.$slots, "title", {}, () => {
        _push(`${ssrInterpolate(__props.title)}`);
      }, _push, _parent);
      _push(`</h1>`);
      if (((_a = __props.links) == null ? void 0 : _a.length) || _ctx.$slots.links) {
        _push(`<div class="${ssrRenderClass(unref(ui).links)}">`);
        ssrRenderSlot(_ctx.$slots, "links", {}, () => {
          _push(`<!--[-->`);
          ssrRenderList(__props.links, (link, index) => {
            _push(ssrRenderComponent(_component_UButton, mergeProps({
              key: index,
              ref_for: true
            }, { ...link, target: link.target || "_blank", color: link.color || "white" }, {
              onClick: link.click
            }), null, _parent));
          });
          _push(`<!--]-->`);
        }, _push, _parent);
        _push(`</div>`);
      } else {
        _push(`<!---->`);
      }
      _push(`</div>`);
      if (__props.description || _ctx.$slots.description) {
        _push(`<div class="${ssrRenderClass(unref(ui).description)}">`);
        ssrRenderSlot(_ctx.$slots, "description", {}, () => {
          _push(`${ssrInterpolate(__props.description)}`);
        }, _push, _parent);
        _push(`</div>`);
      } else {
        _push(`<!---->`);
      }
      ssrRenderSlot(_ctx.$slots, "default", {}, null, _push, _parent);
      _push(`</div></div></div>`);
    };
  }
});
const _sfc_setup = _sfc_main.setup;
_sfc_main.setup = (props, ctx) => {
  const ssrContext = useSSRContext();
  (ssrContext.modules || (ssrContext.modules = /* @__PURE__ */ new Set())).add("node_modules/@nuxt/ui-pro/components/page/PageHeader.vue");
  return _sfc_setup ? _sfc_setup(props, ctx) : void 0;
};
export {
  _sfc_main as _
};
//# sourceMappingURL=PageHeader-maaT1_jD.js.map
