import { a as __nuxt_component_2, u as useUI, e as getULinkProps, _ as __nuxt_component_3, f as __nuxt_component_1$2 } from './Button-CDYqYkPm.mjs';
import { useSSRContext, defineComponent, toRef, computed, ref, watch, withCtx, createVNode, renderSlot, mergeProps, unref, provide, inject, openBlock, createBlock, createCommentVNode, toDisplayString, createSlots, renderList } from 'vue';
import { ssrRenderComponent, ssrRenderSlot, ssrRenderAttrs, ssrRenderClass, ssrInterpolate, ssrRenderList } from 'vue/server-renderer';
import { twMerge, twJoin } from 'tailwind-merge';
import { b as useId, _ as __nuxt_component_1$1 } from './keyboard-D9T68tnc.mjs';
import { q as klona, F as parse, G as getRequestHeader, H as destr, B as isEqual, I as setCookie, J as getCookie, K as deleteCookie } from '../runtime.mjs';
import { m as mergeConfig, b as appConfig, a as useRoute, c as useAppConfig, d as useNuxtApp } from './server.mjs';
import { createSharedComposable, useBreakpoints, breakpointsTailwind, useStorage } from '@vueuse/core';
import { _ as __nuxt_component_2$1 } from './Badge-OGRmTTcV.mjs';
import { _ as __nuxt_component_1 } from './Avatar-GrmG4wkZ.mjs';
import { _ as _export_sfc } from './_plugin-vue_export-helper-1tPrXgE0.mjs';
import { _ as __nuxt_component_0 } from './Accordion-QAMf0_up.mjs';
import './nuxt-link-CgTv0-CK.mjs';
import './Icon-Dm7V5HYH.mjs';
import '@iconify/vue/dist/offline';
import '@iconify/vue';
import './index-BETvu_7i.mjs';
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
import '../routes/renderer.mjs';
import 'vue-bundle-renderer/runtime';
import 'devalue';
import '@unhead/ssr';
import 'unhead';
import '@unhead/shared';
import 'vue-router';

const divider = {
  wrapper: {
    base: "flex items-center align-center text-center",
    horizontal: "w-full flex-row",
    vertical: "flex-col"
  },
  container: {
    base: "font-medium text-gray-700 dark:text-gray-200 flex",
    horizontal: "mx-3 whitespace-nowrap",
    vertical: "my-2"
  },
  border: {
    base: "flex border-gray-200 dark:border-gray-800",
    horizontal: "w-full",
    vertical: "h-full",
    size: {
      horizontal: {
        "2xs": "border-t",
        xs: "border-t-[2px]",
        sm: "border-t-[3px]",
        md: "border-t-[4px]",
        lg: "border-t-[5px]",
        xl: "border-t-[6px]"
      },
      vertical: {
        "2xs": "border-s",
        xs: "border-s-[2px]",
        sm: "border-s-[3px]",
        md: "border-s-[4px]",
        lg: "border-s-[5px]",
        xl: "border-s-[6px]"
      }
    },
    type: {
      solid: "border-solid",
      dotted: "border-dotted",
      dashed: "border-dashed"
    }
  },
  icon: {
    base: "flex-shrink-0 w-5 h-5"
  },
  avatar: {
    base: "flex-shrink-0",
    size: "2xs"
  },
  label: "text-sm",
  default: {
    size: "2xs"
  }
};
const _sfc_main$c = /* @__PURE__ */ defineComponent({
  ...{
    inheritAttrs: false
  },
  __name: "DashboardLayout",
  __ssrInlineRender: true,
  props: {
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
    const config2 = {
      wrapper: "fixed inset-0 flex overflow-hidden"
    };
    const props = __props;
    const { ui, attrs } = useUI("dashboard.layout", toRef(props, "ui"), config2, toRef(props, "class"), true);
    return (_ctx, _push, _parent, _attrs) => {
      _push(`<div${ssrRenderAttrs(mergeProps({
        class: unref(ui).wrapper
      }, unref(attrs), _attrs))}>`);
      ssrRenderSlot(_ctx.$slots, "default", {}, null, _push, _parent);
      _push(`</div>`);
    };
  }
});
const _sfc_setup$c = _sfc_main$c.setup;
_sfc_main$c.setup = (props, ctx) => {
  const ssrContext = useSSRContext();
  (ssrContext.modules || (ssrContext.modules = /* @__PURE__ */ new Set())).add("node_modules/@nuxt/ui-pro/components/dashboard/DashboardLayout.vue");
  return _sfc_setup$c ? _sfc_setup$c(props, ctx) : void 0;
};
const _sfc_main$b = /* @__PURE__ */ defineComponent({
  ...{
    inheritAttrs: false
  },
  __name: "DashboardPanelHandle",
  __ssrInlineRender: true,
  props: {
    orientation: {
      type: String,
      default: "vertical"
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
    const props = __props;
    const config2 = computed(() => {
      const wrapper = twJoin(
        "hidden md:block bg-transparent select-none absolute z-50 group",
        props.orientation === "vertical" && "w-[9px] h-full inset-y-0 -right-[5px] cursor-col-resize",
        props.orientation === "horizontal" && "h-[9px] w-full inset-x-0 -top-[5px] cursor-row-resize"
      );
      const container = twJoin(
        "group-hover:bg-gray-300 dark:group-hover:bg-gray-700 transition duration-200 absolute",
        props.orientation === "vertical" && "w-px h-full inset-x-0 mx-auto",
        props.orientation === "horizontal" && "h-px w-full inset-y-0 my-auto"
      );
      return {
        wrapper,
        container
      };
    });
    const { ui, attrs } = useUI("dashboard.panel.handle", toRef(props, "ui"), config2, toRef(props, "class"), true);
    return (_ctx, _push, _parent, _attrs) => {
      _push(`<div${ssrRenderAttrs(mergeProps(unref(attrs), {
        class: unref(ui).wrapper
      }, _attrs))}><div class="${ssrRenderClass(unref(ui).container)}"></div></div>`);
    };
  }
});
const _sfc_setup$b = _sfc_main$b.setup;
_sfc_main$b.setup = (props, ctx) => {
  const ssrContext = useSSRContext();
  (ssrContext.modules || (ssrContext.modules = /* @__PURE__ */ new Set())).add("node_modules/@nuxt/ui-pro/components/dashboard/DashboardPanelHandle.vue");
  return _sfc_setup$b ? _sfc_setup$b(props, ctx) : void 0;
};
function useRequestEvent(nuxtApp = useNuxtApp()) {
  var _a;
  return (_a = nuxtApp.ssrContext) == null ? void 0 : _a.event;
}
const CookieDefaults = {
  path: "/",
  watch: true,
  decode: (val) => destr(decodeURIComponent(val)),
  encode: (val) => encodeURIComponent(typeof val === "string" ? val : JSON.stringify(val))
};
function useCookie(name, _opts) {
  var _a2;
  var _a;
  const opts = { ...CookieDefaults, ..._opts };
  const cookies = readRawCookies(opts) || {};
  let delay;
  if (opts.maxAge !== void 0) {
    delay = opts.maxAge * 1e3;
  } else if (opts.expires) {
    delay = opts.expires.getTime() - Date.now();
  }
  const hasExpired = delay !== void 0 && delay <= 0;
  const cookieValue = klona(hasExpired ? void 0 : (_a2 = cookies[name]) != null ? _a2 : (_a = opts.default) == null ? void 0 : _a.call(opts));
  const cookie = ref(cookieValue);
  {
    const nuxtApp = useNuxtApp();
    const writeFinalCookieValue = () => {
      if (opts.readonly || isEqual(cookie.value, cookies[name])) {
        return;
      }
      writeServerCookie(useRequestEvent(nuxtApp), name, cookie.value, opts);
    };
    const unhook = nuxtApp.hooks.hookOnce("app:rendered", writeFinalCookieValue);
    nuxtApp.hooks.hookOnce("app:error", () => {
      unhook();
      return writeFinalCookieValue();
    });
  }
  return cookie;
}
function readRawCookies(opts = {}) {
  {
    return parse(getRequestHeader(useRequestEvent(), "cookie") || "", opts);
  }
}
function writeServerCookie(event, name, value, opts = {}) {
  if (event) {
    if (value !== null && value !== void 0) {
      return setCookie(event, name, value, opts);
    }
    if (getCookie(event, name) !== void 0) {
      return deleteCookie(event, name, opts);
    }
  }
}
const useResizable = (key, { min, max, value = 0, storage = "cookie" }) => {
  const el = ref(null);
  const width = storage === "cookie" ? useCookie(key, { default: () => value }) : useStorage(key, () => value);
  const isDragging = ref(false);
  function onMouseMove(e, x) {
    let w = el.value.offsetWidth + e.clientX - x;
    if (min) {
      w = Math.max(w, min);
    }
    if (max) {
      w = Math.min(w, max);
    }
    width.value = w;
    return e.clientX;
  }
  function onDrag(e) {
    if (!el.value)
      return;
    let x = e.clientX;
    (void 0).onmousemove = (e2) => {
      isDragging.value = true;
      x = onMouseMove(e2, x);
    };
    (void 0).onmouseup = () => {
      isDragging.value = false;
      (void 0).onmousemove = (void 0).onmouseup = null;
    };
  }
  return {
    el,
    width,
    isDragging,
    onDrag
  };
};
const _useUIState = () => {
  const route = useRoute();
  const isHeaderDialogOpen = ref(false);
  const isContentSearchModalOpen = ref(false);
  const isDashboardSidebarSlideoverOpen = ref(false);
  const isDashboardSearchModalOpen = ref(false);
  function toggleContentSearch() {
    if (isHeaderDialogOpen.value) {
      isHeaderDialogOpen.value = false;
      setTimeout(() => {
        isContentSearchModalOpen.value = !isContentSearchModalOpen.value;
      }, 0);
      return;
    }
    isContentSearchModalOpen.value = !isContentSearchModalOpen.value;
  }
  function toggleDashboardSearch() {
    if (isDashboardSidebarSlideoverOpen.value) {
      isDashboardSidebarSlideoverOpen.value = false;
      setTimeout(() => {
        isDashboardSearchModalOpen.value = !isDashboardSearchModalOpen.value;
      }, 200);
      return;
    }
    isDashboardSearchModalOpen.value = !isDashboardSearchModalOpen.value;
  }
  watch(() => route.path, () => {
    isDashboardSidebarSlideoverOpen.value = false;
  });
  return {
    isHeaderDialogOpen,
    isContentSearchModalOpen,
    /**
     * @deprecated Use the new {@link isDashboardSidebarSlideoverOpen} ref instead.
     */
    isDashboardSidebarSlidoverOpen: isDashboardSidebarSlideoverOpen,
    isDashboardSidebarSlideoverOpen,
    isDashboardSearchModalOpen,
    toggleContentSearch,
    toggleDashboardSearch
  };
};
const useUIState = createSharedComposable(_useUIState);
const _sfc_main$a = /* @__PURE__ */ defineComponent({
  ...{
    inheritAttrs: false
  },
  __name: "DashboardPanel",
  __ssrInlineRender: true,
  props: {
    id: {
      type: String,
      default: void 0
    },
    modelValue: {
      type: Boolean,
      default: void 0
    },
    collapsible: {
      type: Boolean,
      default: false
    },
    side: {
      type: String,
      default: "left"
    },
    grow: {
      type: Boolean,
      default: false
    },
    resizable: {
      // FIXME: This breaks typecheck
      // type: [Boolean, Object] as PropType<boolean | {
      //   min?: number,
      //   max?: number,
      //   value?: number,
      //   storage?: 'cookie' | 'local'
      // }>,
      type: [Boolean, Object],
      default: false
    },
    width: {
      type: Number,
      default: void 0
    },
    breakpoint: {
      type: String,
      default: "lg"
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
  emits: ["update:modelValue"],
  setup(__props, { expose: __expose, emit: __emit }) {
    const config2 = {
      wrapper: "flex-col items-stretch relative w-full",
      border: "border-b lg:border-b-0 lg:border-r border-gray-200 dark:border-gray-800 lg:w-[--width] flex-shrink-0",
      grow: "flex-1",
      collapsible: "hidden lg:flex",
      slideover: "lg:hidden"
    };
    const props = __props;
    const emit = __emit;
    const id = props.id ? `dashboard:panel:${props.id}` : useId("$dashboard:panel");
    const { ui, attrs } = useUI("dashboard.panel", toRef(props, "ui"), config2, toRef(props, "class"), true);
    const { el, width, onDrag, isDragging } = props.resizable ? useResizable(id, { ...typeof props.resizable === "object" ? props.resizable : {}, value: props.width }) : { el: void 0, width: toRef(props.width), onDrag: void 0, isDragging: void 0 };
    const breakpoints = useBreakpoints(breakpointsTailwind);
    const { isDashboardSidebarSlideoverOpen } = useUIState();
    breakpoints.smaller(props.breakpoint);
    const isOpen = computed({
      get() {
        return props.modelValue !== void 0 ? props.modelValue : isDashboardSidebarSlideoverOpen.value;
      },
      set(value) {
        props.modelValue !== void 0 ? emit("update:modelValue", value) : isDashboardSidebarSlideoverOpen.value = value;
      }
    });
    __expose({
      width,
      isDragging
    });
    provide("isOpen", isOpen);
    return (_ctx, _push, _parent, _attrs) => {
      const _component_UDashboardPanelHandle = _sfc_main$b;
      const _component_ClientOnly = __nuxt_component_1$1;
      _push(`<!--[--><div${ssrRenderAttrs(mergeProps({
        ref_key: "el",
        ref: el
      }, { ...unref(attrs), ..._ctx.$attrs }, {
        class: [unref(ui).wrapper, __props.grow ? unref(ui).grow : unref(ui).border, __props.collapsible ? unref(ui).collapsible : "flex"],
        style: { "--width": unref(width) && !__props.grow ? `${unref(width)}px` : void 0 }
      }))}>`);
      ssrRenderSlot(_ctx.$slots, "default", {}, null, _push, _parent);
      ssrRenderSlot(_ctx.$slots, "handle", { onDrag: unref(onDrag) }, () => {
        if (__props.resizable && !__props.grow) {
          _push(ssrRenderComponent(_component_UDashboardPanelHandle, { onMousedown: unref(onDrag) }, null, _parent));
        } else {
          _push(`<!---->`);
        }
      }, _push, _parent);
      _push(`</div>`);
      _push(ssrRenderComponent(_component_ClientOnly, null, {}, _parent));
      _push(`<!--]-->`);
    };
  }
});
const _sfc_setup$a = _sfc_main$a.setup;
_sfc_main$a.setup = (props, ctx) => {
  const ssrContext = useSSRContext();
  (ssrContext.modules || (ssrContext.modules = /* @__PURE__ */ new Set())).add("node_modules/@nuxt/ui-pro/components/dashboard/DashboardPanel.vue");
  return _sfc_setup$a ? _sfc_setup$a(props, ctx) : void 0;
};
const _sfc_main$9 = /* @__PURE__ */ defineComponent({
  __name: "DashboardNavbarToggle",
  __ssrInlineRender: true,
  setup(__props) {
    const appConfig2 = useAppConfig();
    const isOpen = inject("isOpen", void 0);
    return (_ctx, _push, _parent, _attrs) => {
      var _a, _b;
      const _component_UButton = __nuxt_component_3;
      if (unref(isOpen) !== void 0) {
        _push(ssrRenderComponent(_component_UButton, mergeProps({
          icon: unref(appConfig2).ui.icons.menu
        }, (_b = (_a = _ctx.$ui) == null ? void 0 : _a.button) == null ? void 0 : _b.secondary, {
          "aria-label": `${unref(isOpen) ? "Close" : "Open"} sidebar`,
          class: "lg:hidden",
          onClick: ($event) => isOpen.value = !unref(isOpen)
        }, _attrs), null, _parent));
      } else {
        _push(`<!---->`);
      }
    };
  }
});
const _sfc_setup$9 = _sfc_main$9.setup;
_sfc_main$9.setup = (props, ctx) => {
  const ssrContext = useSSRContext();
  (ssrContext.modules || (ssrContext.modules = /* @__PURE__ */ new Set())).add("node_modules/@nuxt/ui-pro/components/dashboard/DashboardNavbarToggle.vue");
  return _sfc_setup$9 ? _sfc_setup$9(props, ctx) : void 0;
};
const _sfc_main$8 = /* @__PURE__ */ defineComponent({
  ...{
    inheritAttrs: false
  },
  __name: "DashboardNavbar",
  __ssrInlineRender: true,
  props: {
    title: {
      type: String,
      default: void 0
    },
    badge: {
      type: [String, Number, Object],
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
    const config2 = {
      wrapper: "h-[--header-height] flex-shrink-0 flex items-center border-b border-gray-200 dark:border-gray-800 px-4 gap-x-4 min-w-0",
      container: "flex items-center justify-between flex-1 gap-x-1.5 min-w-0",
      left: "flex items-stretch gap-1.5 min-w-0",
      title: "flex items-center gap-1.5 font-semibold text-gray-900 dark:text-white min-w-0",
      badge: {
        wrapper: "inline-flex items-center",
        base: "",
        size: "xs",
        color: "primary",
        variant: "subtle"
      },
      center: "hidden lg:flex",
      right: "flex items-stretch flex-shrink-0 gap-1.5"
    };
    const props = __props;
    const { ui, attrs } = useUI("dashboard.navbar", toRef(props, "ui"), config2, toRef(props, "class"), true);
    return (_ctx, _push, _parent, _attrs) => {
      const _component_UDashboardNavbarToggle = _sfc_main$9;
      const _component_UBadge = __nuxt_component_2$1;
      _push(`<div${ssrRenderAttrs(mergeProps({
        class: unref(ui).wrapper
      }, unref(attrs), _attrs))}><div class="${ssrRenderClass(unref(ui).container)}"><div class="${ssrRenderClass(unref(ui).left)}">`);
      ssrRenderSlot(_ctx.$slots, "toggle", {}, () => {
        _push(ssrRenderComponent(_component_UDashboardNavbarToggle, null, null, _parent));
      }, _push, _parent);
      ssrRenderSlot(_ctx.$slots, "left", {}, () => {
        if (__props.title || _ctx.$slots.title) {
          _push(`<h1 class="${ssrRenderClass(unref(ui).title)}">`);
          ssrRenderSlot(_ctx.$slots, "title", {}, () => {
            _push(`<span class="truncate">${ssrInterpolate(__props.title)}</span>`);
          }, _push, _parent);
          _push(`</h1>`);
        } else {
          _push(`<!---->`);
        }
        if (__props.badge || _ctx.$slots.badge) {
          _push(`<div class="${ssrRenderClass(unref(ui).badge.wrapper)}">`);
          ssrRenderSlot(_ctx.$slots, "badge", {}, () => {
            if (__props.badge) {
              _push(ssrRenderComponent(_component_UBadge, mergeProps({
                size: unref(ui).badge.size,
                color: unref(ui).badge.color,
                variant: unref(ui).badge.variant,
                ...typeof __props.badge === "string" || typeof __props.badge === "number" ? { label: __props.badge } : __props.badge
              }, {
                class: unref(ui).badge.base
              }), null, _parent));
            } else {
              _push(`<!---->`);
            }
          }, _push, _parent);
          _push(`</div>`);
        } else {
          _push(`<!---->`);
        }
      }, _push, _parent);
      _push(`</div>`);
      if (_ctx.$slots.center) {
        _push(`<div class="${ssrRenderClass(unref(ui).center)}">`);
        ssrRenderSlot(_ctx.$slots, "center", {}, null, _push, _parent);
        _push(`</div>`);
      } else {
        _push(`<!---->`);
      }
      if (_ctx.$slots.right || _ctx.$slots.center) {
        _push(`<div class="${ssrRenderClass(unref(ui).right)}">`);
        ssrRenderSlot(_ctx.$slots, "right", {}, null, _push, _parent);
        _push(`</div>`);
      } else {
        _push(`<!---->`);
      }
      _push(`</div></div>`);
    };
  }
});
const _sfc_setup$8 = _sfc_main$8.setup;
_sfc_main$8.setup = (props, ctx) => {
  const ssrContext = useSSRContext();
  (ssrContext.modules || (ssrContext.modules = /* @__PURE__ */ new Set())).add("node_modules/@nuxt/ui-pro/components/dashboard/DashboardNavbar.vue");
  return _sfc_setup$8 ? _sfc_setup$8(props, ctx) : void 0;
};
const _sfc_main$7 = /* @__PURE__ */ defineComponent({
  ...{
    inheritAttrs: false
  },
  __name: "DashboardSidebar",
  __ssrInlineRender: true,
  props: {
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
    const config2 = {
      wrapper: "flex flex-col w-full flex-1 relative overflow-hidden",
      container: "flex-grow flex flex-col min-h-0 gap-y-2 py-2",
      header: "w-full flex flex-col px-4",
      body: "flex-1 px-4 flex flex-col gap-y-2 overflow-y-auto",
      footer: "flex items-center justify-between gap-x-1.5 flex-shrink-0 px-4"
    };
    const props = __props;
    const { ui, attrs } = useUI("dashboard.sidebar", toRef(props, "ui"), config2, toRef(props, "class"), true);
    return (_ctx, _push, _parent, _attrs) => {
      _push(`<div${ssrRenderAttrs(mergeProps({
        class: unref(ui).wrapper
      }, unref(attrs), _attrs))}><div class="${ssrRenderClass(unref(ui).container)}">`);
      if (_ctx.$slots.header) {
        _push(`<div class="${ssrRenderClass(unref(ui).header)}">`);
        ssrRenderSlot(_ctx.$slots, "header", {}, null, _push, _parent);
        _push(`</div>`);
      } else {
        _push(`<!---->`);
      }
      _push(`<div class="${ssrRenderClass(unref(ui).body)}">`);
      ssrRenderSlot(_ctx.$slots, "default", {}, null, _push, _parent);
      _push(`</div>`);
      if (_ctx.$slots.footer) {
        _push(`<div class="${ssrRenderClass(unref(ui).footer)}">`);
        ssrRenderSlot(_ctx.$slots, "footer", {}, null, _push, _parent);
        _push(`</div>`);
      } else {
        _push(`<!---->`);
      }
      _push(`</div></div>`);
    };
  }
});
const _sfc_setup$7 = _sfc_main$7.setup;
_sfc_main$7.setup = (props, ctx) => {
  const ssrContext = useSSRContext();
  (ssrContext.modules || (ssrContext.modules = /* @__PURE__ */ new Set())).add("node_modules/@nuxt/ui-pro/components/dashboard/DashboardSidebar.vue");
  return _sfc_setup$7 ? _sfc_setup$7(props, ctx) : void 0;
};
const _sfc_main$6 = /* @__PURE__ */ defineComponent({
  ...{
    inheritAttrs: false
  },
  __name: "NavigationLinks",
  __ssrInlineRender: true,
  props: {
    level: {
      type: Number,
      default: 0
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
    const config2 = {
      wrapper: "space-y-3",
      wrapperLevel: "space-y-1.5",
      base: "flex items-center gap-1.5 group",
      active: "text-primary font-medium border-current",
      inactive: "text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200 border-transparent hover:border-gray-500 dark:hover:border-gray-400",
      level: "border-l -ml-px pl-4",
      icon: {
        base: "w-5 h-5 flex-shrink-0"
      },
      badge: {
        base: "rounded-full"
      },
      label: "text-sm/6 truncate"
    };
    const props = __props;
    const { ui, attrs } = useUI("navigation.links", toRef(props, "ui"), config2, toRef(props, "class"), true);
    return (_ctx, _push, _parent, _attrs) => {
      var _a;
      const _component_ULink = __nuxt_component_1$2;
      const _component_UIcon = __nuxt_component_2;
      const _component_UBadge = __nuxt_component_2$1;
      if ((_a = __props.links) == null ? void 0 : _a.length) {
        _push(`<div${ssrRenderAttrs(mergeProps({
          class: __props.level > 0 ? unref(ui).wrapperLevel : unref(ui).wrapper
        }, unref(attrs), _attrs))}><!--[-->`);
        ssrRenderList(__props.links, (link, index) => {
          _push(ssrRenderComponent(_component_ULink, mergeProps({
            key: index,
            ref_for: true
          }, unref(getULinkProps)(link), {
            class: [unref(ui).base, __props.level > 0 && unref(ui).level],
            "active-class": unref(ui).active,
            "inactive-class": unref(ui).inactive,
            onClick: link.click
          }), {
            default: withCtx((_, _push2, _parent2, _scopeId) => {
              if (_push2) {
                if (link.icon) {
                  _push2(ssrRenderComponent(_component_UIcon, {
                    name: link.icon,
                    class: unref(twMerge)(unref(ui).icon.base, link.iconClass)
                  }, null, _parent2, _scopeId));
                } else {
                  _push2(`<!---->`);
                }
                _push2(`<span class="${ssrRenderClass(unref(ui).label)}"${_scopeId}>${ssrInterpolate(link.label)}</span>`);
                ssrRenderSlot(_ctx.$slots, "badge", { link }, () => {
                  if (link.badge) {
                    _push2(ssrRenderComponent(_component_UBadge, mergeProps({ ref_for: true }, typeof link.badge === "string" ? { size: "xs", variant: "subtle", label: link.badge } : { size: "xs", variant: "subtle", ...link.badge }, {
                      class: unref(ui).badge.base
                    }), null, _parent2, _scopeId));
                  } else {
                    _push2(`<!---->`);
                  }
                }, _push2, _parent2, _scopeId);
              } else {
                return [
                  link.icon ? (openBlock(), createBlock(_component_UIcon, {
                    key: 0,
                    name: link.icon,
                    class: unref(twMerge)(unref(ui).icon.base, link.iconClass)
                  }, null, 8, ["name", "class"])) : createCommentVNode("", true),
                  createVNode("span", {
                    class: unref(ui).label
                  }, toDisplayString(link.label), 3),
                  renderSlot(_ctx.$slots, "badge", { link }, () => [
                    link.badge ? (openBlock(), createBlock(_component_UBadge, mergeProps({
                      key: 0,
                      ref_for: true
                    }, typeof link.badge === "string" ? { size: "xs", variant: "subtle", label: link.badge } : { size: "xs", variant: "subtle", ...link.badge }, {
                      class: unref(ui).badge.base
                    }), null, 16, ["class"])) : createCommentVNode("", true)
                  ])
                ];
              }
            }),
            _: 2
          }, _parent));
        });
        _push(`<!--]--></div>`);
      } else {
        _push(`<!---->`);
      }
    };
  }
});
const _sfc_setup$6 = _sfc_main$6.setup;
_sfc_main$6.setup = (props, ctx) => {
  const ssrContext = useSSRContext();
  (ssrContext.modules || (ssrContext.modules = /* @__PURE__ */ new Set())).add("node_modules/@nuxt/ui-pro/components/navigation/NavigationLinks.vue");
  return _sfc_setup$6 ? _sfc_setup$6(props, ctx) : void 0;
};
const config = mergeConfig(appConfig.ui.strategy, appConfig.ui.divider, divider);
const _sfc_main$5 = defineComponent({
  components: {
    UIcon: __nuxt_component_2,
    UAvatar: __nuxt_component_1
  },
  inheritAttrs: false,
  props: {
    label: {
      type: String,
      default: null
    },
    icon: {
      type: String,
      default: null
    },
    avatar: {
      type: Object,
      default: null
    },
    size: {
      type: String,
      default: () => config.default.size,
      validator(value) {
        return Object.keys(config.border.size.horizontal).includes(value) || Object.keys(config.border.size.vertical).includes(value);
      }
    },
    orientation: {
      type: String,
      default: "horizontal",
      validator: (value) => ["horizontal", "vertical"].includes(value)
    },
    type: {
      type: String,
      default: "solid",
      validator: (value) => ["solid", "dotted", "dashed"].includes(value)
    },
    class: {
      type: [String, Object, Array],
      default: () => ""
    },
    ui: {
      type: Object,
      default: () => ({})
    }
  },
  setup(props) {
    const { ui, attrs } = useUI("divider", toRef(props, "ui"), config);
    const wrapperClass = computed(() => {
      return twMerge(twJoin(
        ui.value.wrapper.base,
        ui.value.wrapper[props.orientation]
      ), props.class);
    });
    const containerClass = computed(() => {
      return twJoin(
        ui.value.container.base,
        ui.value.container[props.orientation]
      );
    });
    const borderClass = computed(() => {
      return twJoin(
        ui.value.border.base,
        ui.value.border[props.orientation],
        ui.value.border.size[props.orientation][props.size],
        ui.value.border.type[props.type]
      );
    });
    return {
      // eslint-disable-next-line vue/no-dupe-keys
      ui,
      attrs,
      wrapperClass,
      containerClass,
      borderClass
    };
  }
});
function _sfc_ssrRender(_ctx, _push, _parent, _attrs, $props, $setup, $data, $options) {
  const _component_UIcon = __nuxt_component_2;
  const _component_UAvatar = __nuxt_component_1;
  _push(`<div${ssrRenderAttrs(mergeProps({ class: _ctx.wrapperClass }, _ctx.attrs, _attrs))}><div class="${ssrRenderClass(_ctx.borderClass)}"></div>`);
  if (_ctx.label || _ctx.icon || _ctx.avatar || _ctx.$slots.default) {
    _push(`<!--[--><div class="${ssrRenderClass(_ctx.containerClass)}">`);
    ssrRenderSlot(_ctx.$slots, "default", {}, () => {
      if (_ctx.label) {
        _push(`<span class="${ssrRenderClass(_ctx.ui.label)}">${ssrInterpolate(_ctx.label)}</span>`);
      } else if (_ctx.icon) {
        _push(ssrRenderComponent(_component_UIcon, {
          name: _ctx.icon,
          class: _ctx.ui.icon.base
        }, null, _parent));
      } else if (_ctx.avatar) {
        _push(ssrRenderComponent(_component_UAvatar, mergeProps({ size: _ctx.ui.avatar.size, ..._ctx.avatar }, {
          class: _ctx.ui.avatar.base
        }), null, _parent));
      } else {
        _push(`<!---->`);
      }
    }, _push, _parent);
    _push(`</div><div class="${ssrRenderClass(_ctx.borderClass)}"></div><!--]-->`);
  } else {
    _push(`<!---->`);
  }
  _push(`</div>`);
}
const _sfc_setup$5 = _sfc_main$5.setup;
_sfc_main$5.setup = (props, ctx) => {
  const ssrContext = useSSRContext();
  (ssrContext.modules || (ssrContext.modules = /* @__PURE__ */ new Set())).add("node_modules/@nuxt/ui/dist/runtime/components/layout/Divider.vue");
  return _sfc_setup$5 ? _sfc_setup$5(props, ctx) : void 0;
};
const __nuxt_component_6 = /* @__PURE__ */ _export_sfc(_sfc_main$5, [["ssrRender", _sfc_ssrRender]]);
const _sfc_main$4 = /* @__PURE__ */ defineComponent({
  ...{
    inheritAttrs: false
  },
  __name: "NavigationTree",
  __ssrInlineRender: true,
  props: {
    level: {
      type: Number,
      default: 0
    },
    links: {
      type: Array,
      default: () => []
    },
    multiple: {
      type: Boolean,
      default: true
    },
    defaultOpen: {
      type: [Boolean, Number],
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
    const config2 = {
      wrapper: "space-y-3",
      accordion: {},
      links: {}
    };
    const props = __props;
    const { ui, attrs } = useUI("navigation.tree", toRef(props, "ui"), config2, toRef(props, "class"), true);
    const groups = computed(() => {
      var _a;
      const groups2 = [];
      let group = { type: void 0, children: [] };
      for (const link of props.links) {
        const type = ((_a = link.children) == null ? void 0 : _a.length) ? "accordion" : "link";
        if (!group.type) {
          group.type = type;
        }
        if (group.type === type) {
          group.children.push(link);
        } else {
          groups2.push(group);
          group = { type, children: [link] };
        }
      }
      if (group.children.length) {
        groups2.push(group);
      }
      return groups2;
    });
    return (_ctx, _push, _parent, _attrs) => {
      var _a;
      const _component_UNavigationAccordion = _sfc_main$3;
      const _component_UNavigationLinks = _sfc_main$6;
      if ((_a = unref(groups)) == null ? void 0 : _a.length) {
        _push(`<nav${ssrRenderAttrs(mergeProps({
          class: unref(ui).wrapper
        }, unref(attrs), _attrs))}><!--[-->`);
        ssrRenderList(unref(groups), (group, index) => {
          _push(`<!--[-->`);
          if (group.type === "accordion") {
            _push(ssrRenderComponent(_component_UNavigationAccordion, {
              links: group.children,
              level: __props.level,
              multiple: __props.multiple,
              "default-open": __props.defaultOpen,
              ui: { ...unref(ui).accordion, links: unref(ui).links }
            }, null, _parent));
          } else {
            _push(ssrRenderComponent(_component_UNavigationLinks, {
              links: group.children,
              level: __props.level,
              ui: unref(ui).links
            }, null, _parent));
          }
          _push(`<!--]-->`);
        });
        _push(`<!--]--></nav>`);
      } else {
        _push(`<!---->`);
      }
    };
  }
});
const _sfc_setup$4 = _sfc_main$4.setup;
_sfc_main$4.setup = (props, ctx) => {
  const ssrContext = useSSRContext();
  (ssrContext.modules || (ssrContext.modules = /* @__PURE__ */ new Set())).add("node_modules/@nuxt/ui-pro/components/navigation/NavigationTree.vue");
  return _sfc_setup$4 ? _sfc_setup$4(props, ctx) : void 0;
};
const _sfc_main$3 = /* @__PURE__ */ defineComponent({
  ...{
    inheritAttrs: false
  },
  __name: "NavigationAccordion",
  __ssrInlineRender: true,
  props: {
    level: {
      type: Number,
      default: 0
    },
    links: {
      type: Array,
      default: () => []
    },
    multiple: {
      type: Boolean,
      default: true
    },
    defaultOpen: {
      type: [Boolean, Number],
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
    const appConfig2 = useAppConfig();
    const config2 = computed(() => {
      const wrapper = twJoin(
        "space-y-3",
        props.level > 0 && "border-l border-gray-200 dark:border-gray-800 -ml-px hover:border-gray-300 dark:hover:border-gray-700"
      );
      const tree = twJoin(
        "border-l border-gray-200 dark:border-gray-800",
        props.level > 0 ? "ml-6" : "ml-2.5"
      );
      return {
        wrapper,
        container: "space-y-3",
        item: {
          padding: "",
          color: "text-inherit dark:text-inherit"
        },
        button: {
          base: "flex items-center gap-1.5 group w-full focus-visible:outline-primary",
          active: "text-primary border-current",
          inactive: "border-transparent",
          level: "border-l -ml-px pl-3.5",
          icon: {
            base: "w-5 h-5 flex-shrink-0"
          },
          trailingIcon: {
            name: appConfig2.ui.icons.chevron,
            base: "w-5 h-5 ms-auto transform transition-transform duration-200 flex-shrink-0 mr-1.5",
            active: "text-gray-700 dark:text-gray-200",
            inactive: "text-gray-500 dark:text-gray-400 group-hover:text-gray-700 dark:group-hover:text-gray-200 -rotate-90"
          },
          label: "text-sm/6 font-semibold truncate"
        },
        tree,
        links: {}
      };
    });
    const props = __props;
    const route = useRoute();
    const { ui, attrs } = useUI("navigation.accordion", toRef(props, "ui"), config2, toRef(props, "class"), true);
    const items = computed(() => {
      var _a;
      return (_a = props.links) == null ? void 0 : _a.map((link) => {
        const defaultOpen = !props.defaultOpen || typeof props.defaultOpen === "number" && props.level < props.defaultOpen || link.to && route.path.startsWith(link.to.toString());
        return {
          label: link.label,
          icon: link.icon,
          slot: link.label.toLowerCase(),
          disabled: link.disabled,
          defaultOpen,
          children: link.children
        };
      });
    });
    return (_ctx, _push, _parent, _attrs) => {
      const _component_UAccordion = __nuxt_component_0;
      const _component_ULink = __nuxt_component_1$2;
      const _component_UIcon = __nuxt_component_2;
      const _component_UNavigationTree = _sfc_main$4;
      _push(ssrRenderComponent(_component_UAccordion, mergeProps({
        key: unref(route).path,
        items: unref(items),
        multiple: __props.multiple,
        ui: unref(ui)
      }, unref(attrs), _attrs), createSlots({
        default: withCtx(({ item, open }, _push2, _parent2, _scopeId) => {
          if (_push2) {
            _push2(ssrRenderComponent(_component_ULink, {
              class: [unref(ui).button.base, __props.level > 0 && unref(ui).button.level],
              "active-class": unref(ui).button.active,
              "inactive-class": unref(ui).button.inactive
            }, {
              default: withCtx((_, _push3, _parent3, _scopeId2) => {
                if (_push3) {
                  if (item.icon) {
                    _push3(ssrRenderComponent(_component_UIcon, {
                      name: item.icon,
                      class: unref(ui).button.icon.base
                    }, null, _parent3, _scopeId2));
                  } else {
                    _push3(`<!---->`);
                  }
                  _push3(`<span class="${ssrRenderClass(unref(ui).button.label)}"${_scopeId2}>${ssrInterpolate(item.label)}</span>`);
                  if (!item.disabled) {
                    _push3(ssrRenderComponent(_component_UIcon, {
                      name: unref(ui).button.trailingIcon.name,
                      class: [unref(ui).button.trailingIcon.base, open ? unref(ui).button.trailingIcon.active : unref(ui).button.trailingIcon.inactive]
                    }, null, _parent3, _scopeId2));
                  } else {
                    _push3(`<!---->`);
                  }
                } else {
                  return [
                    item.icon ? (openBlock(), createBlock(_component_UIcon, {
                      key: 0,
                      name: item.icon,
                      class: unref(ui).button.icon.base
                    }, null, 8, ["name", "class"])) : createCommentVNode("", true),
                    createVNode("span", {
                      class: unref(ui).button.label
                    }, toDisplayString(item.label), 3),
                    !item.disabled ? (openBlock(), createBlock(_component_UIcon, {
                      key: 1,
                      name: unref(ui).button.trailingIcon.name,
                      class: [unref(ui).button.trailingIcon.base, open ? unref(ui).button.trailingIcon.active : unref(ui).button.trailingIcon.inactive]
                    }, null, 8, ["name", "class"])) : createCommentVNode("", true)
                  ];
                }
              }),
              _: 2
            }, _parent2, _scopeId));
          } else {
            return [
              createVNode(_component_ULink, {
                class: [unref(ui).button.base, __props.level > 0 && unref(ui).button.level],
                "active-class": unref(ui).button.active,
                "inactive-class": unref(ui).button.inactive
              }, {
                default: withCtx(() => [
                  item.icon ? (openBlock(), createBlock(_component_UIcon, {
                    key: 0,
                    name: item.icon,
                    class: unref(ui).button.icon.base
                  }, null, 8, ["name", "class"])) : createCommentVNode("", true),
                  createVNode("span", {
                    class: unref(ui).button.label
                  }, toDisplayString(item.label), 3),
                  !item.disabled ? (openBlock(), createBlock(_component_UIcon, {
                    key: 1,
                    name: unref(ui).button.trailingIcon.name,
                    class: [unref(ui).button.trailingIcon.base, open ? unref(ui).button.trailingIcon.active : unref(ui).button.trailingIcon.inactive]
                  }, null, 8, ["name", "class"])) : createCommentVNode("", true)
                ]),
                _: 2
              }, 1032, ["class", "active-class", "inactive-class"])
            ];
          }
        }),
        _: 2
      }, [
        renderList(__props.links, ({ label }, index) => {
          return {
            name: label.toLowerCase(),
            fn: withCtx(({ item }, _push2, _parent2, _scopeId) => {
              if (_push2) {
                _push2(ssrRenderComponent(_component_UNavigationTree, {
                  links: item.children,
                  level: __props.level + 1,
                  "default-open": __props.defaultOpen,
                  multiple: __props.multiple,
                  class: unref(ui).tree,
                  ui: { links: unref(ui).links }
                }, null, _parent2, _scopeId));
              } else {
                return [
                  createVNode(_component_UNavigationTree, {
                    links: item.children,
                    level: __props.level + 1,
                    "default-open": __props.defaultOpen,
                    multiple: __props.multiple,
                    class: unref(ui).tree,
                    ui: { links: unref(ui).links }
                  }, null, 8, ["links", "level", "default-open", "multiple", "class", "ui"])
                ];
              }
            })
          };
        })
      ]), _parent));
    };
  }
});
const _sfc_setup$3 = _sfc_main$3.setup;
_sfc_main$3.setup = (props, ctx) => {
  const ssrContext = useSSRContext();
  (ssrContext.modules || (ssrContext.modules = /* @__PURE__ */ new Set())).add("node_modules/@nuxt/ui-pro/components/navigation/NavigationAccordion.vue");
  return _sfc_setup$3 ? _sfc_setup$3(props, ctx) : void 0;
};
const _sfc_main$2 = /* @__PURE__ */ defineComponent({
  ...{
    inheritAttrs: false
  },
  __name: "DashboardPage",
  __ssrInlineRender: true,
  props: {
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
    const config2 = {
      wrapper: "flex flex-1 w-full min-w-0"
    };
    const props = __props;
    const { ui, attrs } = useUI("dashboard.page", toRef(props, "ui"), config2, toRef(props, "class"), true);
    return (_ctx, _push, _parent, _attrs) => {
      _push(`<div${ssrRenderAttrs(mergeProps({
        class: unref(ui).wrapper
      }, unref(attrs), _attrs))}>`);
      ssrRenderSlot(_ctx.$slots, "default", {}, null, _push, _parent);
      _push(`</div>`);
    };
  }
});
const _sfc_setup$2 = _sfc_main$2.setup;
_sfc_main$2.setup = (props, ctx) => {
  const ssrContext = useSSRContext();
  (ssrContext.modules || (ssrContext.modules = /* @__PURE__ */ new Set())).add("node_modules/@nuxt/ui-pro/components/dashboard/DashboardPage.vue");
  return _sfc_setup$2 ? _sfc_setup$2(props, ctx) : void 0;
};
const _sfc_main$1 = /* @__PURE__ */ defineComponent({
  ...{
    inheritAttrs: false
  },
  __name: "DashboardPanelContent",
  __ssrInlineRender: true,
  props: {
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
    const config2 = {
      wrapper: "p-4 flex-1 flex flex-col overflow-y-auto"
    };
    const props = __props;
    const { ui, attrs } = useUI("dashboard.panel.content", toRef(props, "ui"), config2, toRef(props, "class"), true);
    return (_ctx, _push, _parent, _attrs) => {
      _push(`<div${ssrRenderAttrs(mergeProps({
        class: unref(ui).wrapper
      }, unref(attrs), _attrs))}>`);
      ssrRenderSlot(_ctx.$slots, "default", {}, null, _push, _parent);
      _push(`</div>`);
    };
  }
});
const _sfc_setup$1 = _sfc_main$1.setup;
_sfc_main$1.setup = (props, ctx) => {
  const ssrContext = useSSRContext();
  (ssrContext.modules || (ssrContext.modules = /* @__PURE__ */ new Set())).add("node_modules/@nuxt/ui-pro/components/dashboard/DashboardPanelContent.vue");
  return _sfc_setup$1 ? _sfc_setup$1(props, ctx) : void 0;
};
const _sfc_main = {
  __name: "default",
  __ssrInlineRender: true,
  setup(__props) {
    const containers = ref([]);
    const projects = ref([]);
    const databases = ref([]);
    const containersCollapsed = ref(false);
    const otherContainersCollapsed = ref(false);
    const projectsCollapsed = ref(false);
    const databasesCollapsed = ref(false);
    const saveCollapsedState = (key, newValue) => {
      localStorage.setItem(key, newValue);
    };
    const links = ref([
      {
        id: "dashboard",
        label: "Dashboard",
        icon: "i-heroicons-rectangle-group",
        to: "/"
      }
    ]);
    const containersLinks = ref([
      {
        label: "Containers",
        icon: "i-heroicons-cube",
        children: [
          {
            id: "containers",
            label: "Overview",
            to: "/containers"
          }
        ]
      }
    ]);
    const otherContainersLinks = ref([
      {
        id: "other-containers",
        label: "Other Containers",
        icon: "i-heroicons-cube-transparent",
        children: [
          {
            id: "other-containers",
            label: "Overview",
            to: "/other-containers"
          }
        ]
      }
    ]);
    const projectsLinks = ref([
      {
        label: "Projects",
        icon: "i-heroicons-folder",
        children: [
          {
            id: "projects",
            label: "Overview",
            to: "/projects"
          }
        ]
      }
    ]);
    const databasesLinks = ref([
      {
        label: "Databases",
        icon: "i-heroicons-circle-stack",
        children: [
          {
            id: "databases",
            label: "Overview",
            to: "/databases"
          }
        ]
      }
    ]);
    watch(containers, (newContainers) => {
      containersLinks.value = [
        {
          label: "Containers",
          icon: "i-heroicons-cube",
          children: [
            {
              id: "containers",
              label: "Overview",
              to: "/containers"
            },
            ...newContainers.map((container) => ({
              id: container.getId(),
              label: container.getName(),
              to: `/containers/${container.getId()}`,
              badge: {
                color: container.isRunning() ? "green" : "red",
                label: container.getStatus()
              }
            }))
          ]
        }
      ];
    });
    watch(projects, (newProjects) => {
      projectsLinks.value = [
        {
          id: "projects",
          label: "Projects",
          icon: "i-heroicons-folder",
          children: [
            {
              id: "projects",
              label: "Overview",
              to: "/projects"
            },
            ...newProjects.map((project) => ({
              id: project.getPath(),
              label: project.getName(),
              to: `/projects/${project.getName()}`,
              badge: {
                color: project.getType().color,
                label: project.getType().label
              }
            }))
          ]
        }
      ];
    });
    watch(databases, (newDatabases) => {
      databasesLinks.value = [
        {
          label: "Databases",
          icon: "i-heroicons-circle-stack",
          children: [
            {
              id: "databases",
              label: "Overview",
              to: "/databases"
            },
            ...newDatabases.map((database) => ({
              id: database.getName(),
              label: database.getName(),
              to: `/databases/${database.getContainerName()}/${database.getName()}`,
              badge: {
                color: database.getType().color,
                label: database.getType().label
              }
            }))
          ]
        }
      ];
    });
    watch(containersCollapsed, (newValue) => {
      saveCollapsedState("containers-collapsed", newValue);
    });
    watch(projectsCollapsed, (newValue) => {
      saveCollapsedState("projects-collapsed", newValue);
    });
    watch(databasesCollapsed, (newValue) => {
      saveCollapsedState("databases-collapsed", newValue);
    });
    return (_ctx, _push, _parent, _attrs) => {
      const _component_UDashboardLayout = _sfc_main$c;
      const _component_UDashboardPanel = _sfc_main$a;
      const _component_UDashboardNavbar = _sfc_main$8;
      const _component_UButton = __nuxt_component_3;
      const _component_UDashboardSidebar = _sfc_main$7;
      const _component_UNavigationLinks = _sfc_main$6;
      const _component_UDivider = __nuxt_component_6;
      const _component_UNavigationAccordion = _sfc_main$3;
      const _component_UDashboardPage = _sfc_main$2;
      const _component_UDashboardPanelContent = _sfc_main$1;
      _push(ssrRenderComponent(_component_UDashboardLayout, _attrs, {
        default: withCtx((_, _push2, _parent2, _scopeId) => {
          if (_push2) {
            _push2(ssrRenderComponent(_component_UDashboardPanel, {
              width: 250,
              resizable: { min: 200, max: 300 },
              collapsible: ""
            }, {
              default: withCtx((_2, _push3, _parent3, _scopeId2) => {
                if (_push3) {
                  _push3(ssrRenderComponent(_component_UDashboardNavbar, {
                    class: "!border-transparent",
                    ui: { left: "flex-1" }
                  }, {
                    left: withCtx((_3, _push4, _parent4, _scopeId3) => {
                      if (_push4) {
                        _push4(ssrRenderComponent(_component_UButton, {
                          icon: "i-heroicons-square-3-stack-3d",
                          size: "xl",
                          color: "white",
                          square: "",
                          variant: "ghost",
                          label: "Devner",
                          class: "pl-0 uppercase"
                        }, null, _parent4, _scopeId3));
                      } else {
                        return [
                          createVNode(_component_UButton, {
                            icon: "i-heroicons-square-3-stack-3d",
                            size: "xl",
                            color: "white",
                            square: "",
                            variant: "ghost",
                            label: "Devner",
                            class: "pl-0 uppercase"
                          })
                        ];
                      }
                    }),
                    _: 1
                  }, _parent3, _scopeId2));
                  _push3(ssrRenderComponent(_component_UDashboardSidebar, null, {
                    default: withCtx((_3, _push4, _parent4, _scopeId3) => {
                      if (_push4) {
                        _push4(ssrRenderComponent(_component_UNavigationLinks, { links: links.value }, null, _parent4, _scopeId3));
                        _push4(ssrRenderComponent(_component_UDivider, { class: "my-4" }, null, _parent4, _scopeId3));
                        _push4(ssrRenderComponent(_component_UNavigationAccordion, {
                          links: containersLinks.value,
                          defaultOpen: containersCollapsed.value
                        }, null, _parent4, _scopeId3));
                        _push4(ssrRenderComponent(_component_UNavigationAccordion, {
                          links: projectsLinks.value,
                          defaultOpen: projectsCollapsed.value
                        }, null, _parent4, _scopeId3));
                        _push4(ssrRenderComponent(_component_UNavigationAccordion, {
                          links: databasesLinks.value,
                          defaultOpen: databasesCollapsed.value
                        }, null, _parent4, _scopeId3));
                        _push4(`<div class="flex-1"${_scopeId3}></div>`);
                        _push4(ssrRenderComponent(_component_UNavigationAccordion, {
                          links: otherContainersLinks.value,
                          defaultOpen: otherContainersCollapsed.value
                        }, null, _parent4, _scopeId3));
                      } else {
                        return [
                          createVNode(_component_UNavigationLinks, { links: links.value }, null, 8, ["links"]),
                          createVNode(_component_UDivider, { class: "my-4" }),
                          createVNode(_component_UNavigationAccordion, {
                            links: containersLinks.value,
                            defaultOpen: containersCollapsed.value
                          }, null, 8, ["links", "defaultOpen"]),
                          createVNode(_component_UNavigationAccordion, {
                            links: projectsLinks.value,
                            defaultOpen: projectsCollapsed.value
                          }, null, 8, ["links", "defaultOpen"]),
                          createVNode(_component_UNavigationAccordion, {
                            links: databasesLinks.value,
                            defaultOpen: databasesCollapsed.value
                          }, null, 8, ["links", "defaultOpen"]),
                          createVNode("div", { class: "flex-1" }),
                          createVNode(_component_UNavigationAccordion, {
                            links: otherContainersLinks.value,
                            defaultOpen: otherContainersCollapsed.value
                          }, null, 8, ["links", "defaultOpen"])
                        ];
                      }
                    }),
                    _: 1
                  }, _parent3, _scopeId2));
                } else {
                  return [
                    createVNode(_component_UDashboardNavbar, {
                      class: "!border-transparent",
                      ui: { left: "flex-1" }
                    }, {
                      left: withCtx(() => [
                        createVNode(_component_UButton, {
                          icon: "i-heroicons-square-3-stack-3d",
                          size: "xl",
                          color: "white",
                          square: "",
                          variant: "ghost",
                          label: "Devner",
                          class: "pl-0 uppercase"
                        })
                      ]),
                      _: 1
                    }),
                    createVNode(_component_UDashboardSidebar, null, {
                      default: withCtx(() => [
                        createVNode(_component_UNavigationLinks, { links: links.value }, null, 8, ["links"]),
                        createVNode(_component_UDivider, { class: "my-4" }),
                        createVNode(_component_UNavigationAccordion, {
                          links: containersLinks.value,
                          defaultOpen: containersCollapsed.value
                        }, null, 8, ["links", "defaultOpen"]),
                        createVNode(_component_UNavigationAccordion, {
                          links: projectsLinks.value,
                          defaultOpen: projectsCollapsed.value
                        }, null, 8, ["links", "defaultOpen"]),
                        createVNode(_component_UNavigationAccordion, {
                          links: databasesLinks.value,
                          defaultOpen: databasesCollapsed.value
                        }, null, 8, ["links", "defaultOpen"]),
                        createVNode("div", { class: "flex-1" }),
                        createVNode(_component_UNavigationAccordion, {
                          links: otherContainersLinks.value,
                          defaultOpen: otherContainersCollapsed.value
                        }, null, 8, ["links", "defaultOpen"])
                      ]),
                      _: 1
                    })
                  ];
                }
              }),
              _: 1
            }, _parent2, _scopeId));
            _push2(ssrRenderComponent(_component_UDashboardPage, null, {
              default: withCtx((_2, _push3, _parent3, _scopeId2) => {
                if (_push3) {
                  _push3(ssrRenderComponent(_component_UDashboardPanel, { grow: "" }, {
                    default: withCtx((_3, _push4, _parent4, _scopeId3) => {
                      if (_push4) {
                        _push4(ssrRenderComponent(_component_UDashboardPanelContent, { class: "pt-0" }, {
                          default: withCtx((_4, _push5, _parent5, _scopeId4) => {
                            if (_push5) {
                              ssrRenderSlot(_ctx.$slots, "default", {}, null, _push5, _parent5, _scopeId4);
                            } else {
                              return [
                                renderSlot(_ctx.$slots, "default")
                              ];
                            }
                          }),
                          _: 3
                        }, _parent4, _scopeId3));
                      } else {
                        return [
                          createVNode(_component_UDashboardPanelContent, { class: "pt-0" }, {
                            default: withCtx(() => [
                              renderSlot(_ctx.$slots, "default")
                            ]),
                            _: 3
                          })
                        ];
                      }
                    }),
                    _: 3
                  }, _parent3, _scopeId2));
                } else {
                  return [
                    createVNode(_component_UDashboardPanel, { grow: "" }, {
                      default: withCtx(() => [
                        createVNode(_component_UDashboardPanelContent, { class: "pt-0" }, {
                          default: withCtx(() => [
                            renderSlot(_ctx.$slots, "default")
                          ]),
                          _: 3
                        })
                      ]),
                      _: 3
                    })
                  ];
                }
              }),
              _: 3
            }, _parent2, _scopeId));
          } else {
            return [
              createVNode(_component_UDashboardPanel, {
                width: 250,
                resizable: { min: 200, max: 300 },
                collapsible: ""
              }, {
                default: withCtx(() => [
                  createVNode(_component_UDashboardNavbar, {
                    class: "!border-transparent",
                    ui: { left: "flex-1" }
                  }, {
                    left: withCtx(() => [
                      createVNode(_component_UButton, {
                        icon: "i-heroicons-square-3-stack-3d",
                        size: "xl",
                        color: "white",
                        square: "",
                        variant: "ghost",
                        label: "Devner",
                        class: "pl-0 uppercase"
                      })
                    ]),
                    _: 1
                  }),
                  createVNode(_component_UDashboardSidebar, null, {
                    default: withCtx(() => [
                      createVNode(_component_UNavigationLinks, { links: links.value }, null, 8, ["links"]),
                      createVNode(_component_UDivider, { class: "my-4" }),
                      createVNode(_component_UNavigationAccordion, {
                        links: containersLinks.value,
                        defaultOpen: containersCollapsed.value
                      }, null, 8, ["links", "defaultOpen"]),
                      createVNode(_component_UNavigationAccordion, {
                        links: projectsLinks.value,
                        defaultOpen: projectsCollapsed.value
                      }, null, 8, ["links", "defaultOpen"]),
                      createVNode(_component_UNavigationAccordion, {
                        links: databasesLinks.value,
                        defaultOpen: databasesCollapsed.value
                      }, null, 8, ["links", "defaultOpen"]),
                      createVNode("div", { class: "flex-1" }),
                      createVNode(_component_UNavigationAccordion, {
                        links: otherContainersLinks.value,
                        defaultOpen: otherContainersCollapsed.value
                      }, null, 8, ["links", "defaultOpen"])
                    ]),
                    _: 1
                  })
                ]),
                _: 1
              }),
              createVNode(_component_UDashboardPage, null, {
                default: withCtx(() => [
                  createVNode(_component_UDashboardPanel, { grow: "" }, {
                    default: withCtx(() => [
                      createVNode(_component_UDashboardPanelContent, { class: "pt-0" }, {
                        default: withCtx(() => [
                          renderSlot(_ctx.$slots, "default")
                        ]),
                        _: 3
                      })
                    ]),
                    _: 3
                  })
                ]),
                _: 3
              })
            ];
          }
        }),
        _: 3
      }, _parent));
    };
  }
};
const _sfc_setup = _sfc_main.setup;
_sfc_main.setup = (props, ctx) => {
  const ssrContext = useSSRContext();
  (ssrContext.modules || (ssrContext.modules = /* @__PURE__ */ new Set())).add("layouts/default.vue");
  return _sfc_setup ? _sfc_setup(props, ctx) : void 0;
};

export { _sfc_main as default };
//# sourceMappingURL=default-iW-ycMlB.mjs.map
