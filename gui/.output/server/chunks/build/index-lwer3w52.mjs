import { _ as _sfc_main$1 } from './PageHeader-maaT1_jD.mjs';
import { _ as __nuxt_component_0 } from './Accordion-QAMf0_up.mjs';
import { _ as _sfc_main$1$1, a as __nuxt_component_3, C as Container } from './Container-B3wCur7c.mjs';
import { _ as __nuxt_component_2 } from './Badge-OGRmTTcV.mjs';
import { _ as __nuxt_component_3$1 } from './Button-CDYqYkPm.mjs';
import { ref, watch, withCtx, createTextVNode, toDisplayString, createVNode, openBlock, createBlock, createCommentVNode, Fragment, renderList, useSSRContext } from 'vue';
import { ssrRenderComponent, ssrRenderList, ssrInterpolate } from 'vue/server-renderer';
import './keyboard-D9T68tnc.mjs';
import './server.mjs';
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
import '../routes/renderer.mjs';
import 'vue-bundle-renderer/runtime';
import 'devalue';
import '@unhead/ssr';
import 'unhead';
import '@unhead/shared';
import 'vue-router';
import '@vueuse/core';
import 'tailwind-merge';
import './_plugin-vue_export-helper-1tPrXgE0.mjs';
import 'strip-ansi';
import './nuxt-link-CgTv0-CK.mjs';
import './Icon-Dm7V5HYH.mjs';
import '@iconify/vue/dist/offline';
import '@iconify/vue';
import './index-BETvu_7i.mjs';

const _sfc_main = {
  __name: "index",
  __ssrInlineRender: true,
  setup(__props) {
    const containers = ref([]);
    const loadContainers = async () => {
      try {
        containers.value = await Container.allOthers();
        console.log("Loaded containers:", containers.value);
      } catch (error) {
        console.error("Failed to load containers:", error);
      }
    };
    const performContainerAction = async (container, action) => {
      try {
        await container[action]();
        loadContainers();
      } catch (error) {
        console.error(`Failed to ${action} container:`, error);
      }
    };
    watch(containers, (newContainers) => {
      containers.value = newContainers;
    });
    return (_ctx, _push, _parent, _attrs) => {
      const _component_UPageHeader = _sfc_main$1;
      const _component_UAccordion = __nuxt_component_0;
      const _component_UPageGrid = _sfc_main$1$1;
      const _component_UCard = __nuxt_component_3;
      const _component_UBadge = __nuxt_component_2;
      const _component_UButton = __nuxt_component_3$1;
      _push(`<!--[-->`);
      _push(ssrRenderComponent(_component_UPageHeader, {
        title: "Other Containers",
        class: "sticky top-0 bg-white z-10 dark:bg-gray-800 mb-4 pb-4"
      }, null, _parent));
      if (containers.value.length) {
        _push(ssrRenderComponent(_component_UAccordion, {
          multiple: "",
          color: "white",
          variant: "solid",
          size: "lg",
          items: containers.value
        }, {
          item: withCtx(({ item }, _push2, _parent2, _scopeId) => {
            if (_push2) {
              _push2(ssrRenderComponent(_component_UPageGrid, { class: "py-4" }, {
                default: withCtx((_, _push3, _parent3, _scopeId2) => {
                  if (_push3) {
                    _push3(`<!--[-->`);
                    ssrRenderList(item.containers, (container) => {
                      _push3(ssrRenderComponent(_component_UCard, {
                        key: container.getId()
                      }, {
                        header: withCtx((_2, _push4, _parent4, _scopeId3) => {
                          if (_push4) {
                            _push4(`${ssrInterpolate(container.getName())} - `);
                            _push4(ssrRenderComponent(_component_UBadge, {
                              label: container.getStatus(),
                              color: container.isRunning() ? "green" : "red",
                              variant: "soft"
                            }, null, _parent4, _scopeId3));
                          } else {
                            return [
                              createTextVNode(toDisplayString(container.getName()) + " - ", 1),
                              createVNode(_component_UBadge, {
                                label: container.getStatus(),
                                color: container.isRunning() ? "green" : "red",
                                variant: "soft"
                              }, null, 8, ["label", "color"])
                            ];
                          }
                        }),
                        default: withCtx((_2, _push4, _parent4, _scopeId3) => {
                          if (_push4) {
                            _push4(`<div class="flex flex-wrap gap-2"${_scopeId3}>`);
                            if (!container.isRunning()) {
                              _push4(ssrRenderComponent(_component_UButton, {
                                onClick: () => performContainerAction(container, "start"),
                                icon: container.isLoading("start") ? "i-heroicons-spinner" : "i-heroicons-play",
                                size: "sm",
                                color: "white",
                                square: "",
                                variant: "solid",
                                loading: container.isLoading("start"),
                                label: "Start"
                              }, null, _parent4, _scopeId3));
                            } else {
                              _push4(`<!---->`);
                            }
                            if (container.isRunning()) {
                              _push4(ssrRenderComponent(_component_UButton, {
                                onClick: () => performContainerAction(container, "restart"),
                                icon: container.isLoading("restart") ? "i-heroicons-spinner" : "i-heroicons-arrow-path",
                                size: "sm",
                                color: "white",
                                square: "",
                                variant: "solid",
                                loading: container.isLoading("restart"),
                                label: "Restart"
                              }, null, _parent4, _scopeId3));
                            } else {
                              _push4(`<!---->`);
                            }
                            if (container.isRunning()) {
                              _push4(ssrRenderComponent(_component_UButton, {
                                onClick: () => performContainerAction(container, "stop"),
                                icon: container.isLoading("stop") ? "i-heroicons-spinner" : "i-heroicons-stop",
                                size: "sm",
                                color: "white",
                                square: "",
                                variant: "solid",
                                loading: container.isLoading("stop"),
                                label: "Stop"
                              }, null, _parent4, _scopeId3));
                            } else {
                              _push4(`<!---->`);
                            }
                            _push4(ssrRenderComponent(_component_UButton, {
                              onClick: () => performContainerAction(container, "remove"),
                              icon: container.isLoading("remove") ? "i-heroicons-spinner" : "i-heroicons-trash",
                              size: "sm",
                              color: "white",
                              square: "",
                              variant: "solid",
                              loading: container.isLoading("remove"),
                              label: "Remove"
                            }, null, _parent4, _scopeId3));
                            _push4(ssrRenderComponent(_component_UButton, {
                              onClick: () => performContainerAction(container, "rebuild"),
                              icon: container.isLoading("rebuild") ? "i-heroicons-spinner" : "i-heroicons-arrow-path-rounded-square",
                              size: "sm",
                              color: "white",
                              square: "",
                              variant: "solid",
                              loading: container.isLoading("rebuild"),
                              label: "Rebuild"
                            }, null, _parent4, _scopeId3));
                            _push4(ssrRenderComponent(_component_UButton, {
                              to: `/containers/${container.getId()}`,
                              icon: "i-heroicons-eye",
                              size: "sm",
                              color: "white",
                              square: "",
                              variant: "solid",
                              label: "Details"
                            }, null, _parent4, _scopeId3));
                            _push4(`</div>`);
                          } else {
                            return [
                              createVNode("div", { class: "flex flex-wrap gap-2" }, [
                                !container.isRunning() ? (openBlock(), createBlock(_component_UButton, {
                                  key: 0,
                                  onClick: () => performContainerAction(container, "start"),
                                  icon: container.isLoading("start") ? "i-heroicons-spinner" : "i-heroicons-play",
                                  size: "sm",
                                  color: "white",
                                  square: "",
                                  variant: "solid",
                                  loading: container.isLoading("start"),
                                  label: "Start"
                                }, null, 8, ["onClick", "icon", "loading"])) : createCommentVNode("", true),
                                container.isRunning() ? (openBlock(), createBlock(_component_UButton, {
                                  key: 1,
                                  onClick: () => performContainerAction(container, "restart"),
                                  icon: container.isLoading("restart") ? "i-heroicons-spinner" : "i-heroicons-arrow-path",
                                  size: "sm",
                                  color: "white",
                                  square: "",
                                  variant: "solid",
                                  loading: container.isLoading("restart"),
                                  label: "Restart"
                                }, null, 8, ["onClick", "icon", "loading"])) : createCommentVNode("", true),
                                container.isRunning() ? (openBlock(), createBlock(_component_UButton, {
                                  key: 2,
                                  onClick: () => performContainerAction(container, "stop"),
                                  icon: container.isLoading("stop") ? "i-heroicons-spinner" : "i-heroicons-stop",
                                  size: "sm",
                                  color: "white",
                                  square: "",
                                  variant: "solid",
                                  loading: container.isLoading("stop"),
                                  label: "Stop"
                                }, null, 8, ["onClick", "icon", "loading"])) : createCommentVNode("", true),
                                createVNode(_component_UButton, {
                                  onClick: () => performContainerAction(container, "remove"),
                                  icon: container.isLoading("remove") ? "i-heroicons-spinner" : "i-heroicons-trash",
                                  size: "sm",
                                  color: "white",
                                  square: "",
                                  variant: "solid",
                                  loading: container.isLoading("remove"),
                                  label: "Remove"
                                }, null, 8, ["onClick", "icon", "loading"]),
                                createVNode(_component_UButton, {
                                  onClick: () => performContainerAction(container, "rebuild"),
                                  icon: container.isLoading("rebuild") ? "i-heroicons-spinner" : "i-heroicons-arrow-path-rounded-square",
                                  size: "sm",
                                  color: "white",
                                  square: "",
                                  variant: "solid",
                                  loading: container.isLoading("rebuild"),
                                  label: "Rebuild"
                                }, null, 8, ["onClick", "icon", "loading"]),
                                createVNode(_component_UButton, {
                                  to: `/containers/${container.getId()}`,
                                  icon: "i-heroicons-eye",
                                  size: "sm",
                                  color: "white",
                                  square: "",
                                  variant: "solid",
                                  label: "Details"
                                }, null, 8, ["to"])
                              ])
                            ];
                          }
                        }),
                        _: 2
                      }, _parent3, _scopeId2));
                    });
                    _push3(`<!--]-->`);
                  } else {
                    return [
                      (openBlock(true), createBlock(Fragment, null, renderList(item.containers, (container) => {
                        return openBlock(), createBlock(_component_UCard, {
                          key: container.getId()
                        }, {
                          header: withCtx(() => [
                            createTextVNode(toDisplayString(container.getName()) + " - ", 1),
                            createVNode(_component_UBadge, {
                              label: container.getStatus(),
                              color: container.isRunning() ? "green" : "red",
                              variant: "soft"
                            }, null, 8, ["label", "color"])
                          ]),
                          default: withCtx(() => [
                            createVNode("div", { class: "flex flex-wrap gap-2" }, [
                              !container.isRunning() ? (openBlock(), createBlock(_component_UButton, {
                                key: 0,
                                onClick: () => performContainerAction(container, "start"),
                                icon: container.isLoading("start") ? "i-heroicons-spinner" : "i-heroicons-play",
                                size: "sm",
                                color: "white",
                                square: "",
                                variant: "solid",
                                loading: container.isLoading("start"),
                                label: "Start"
                              }, null, 8, ["onClick", "icon", "loading"])) : createCommentVNode("", true),
                              container.isRunning() ? (openBlock(), createBlock(_component_UButton, {
                                key: 1,
                                onClick: () => performContainerAction(container, "restart"),
                                icon: container.isLoading("restart") ? "i-heroicons-spinner" : "i-heroicons-arrow-path",
                                size: "sm",
                                color: "white",
                                square: "",
                                variant: "solid",
                                loading: container.isLoading("restart"),
                                label: "Restart"
                              }, null, 8, ["onClick", "icon", "loading"])) : createCommentVNode("", true),
                              container.isRunning() ? (openBlock(), createBlock(_component_UButton, {
                                key: 2,
                                onClick: () => performContainerAction(container, "stop"),
                                icon: container.isLoading("stop") ? "i-heroicons-spinner" : "i-heroicons-stop",
                                size: "sm",
                                color: "white",
                                square: "",
                                variant: "solid",
                                loading: container.isLoading("stop"),
                                label: "Stop"
                              }, null, 8, ["onClick", "icon", "loading"])) : createCommentVNode("", true),
                              createVNode(_component_UButton, {
                                onClick: () => performContainerAction(container, "remove"),
                                icon: container.isLoading("remove") ? "i-heroicons-spinner" : "i-heroicons-trash",
                                size: "sm",
                                color: "white",
                                square: "",
                                variant: "solid",
                                loading: container.isLoading("remove"),
                                label: "Remove"
                              }, null, 8, ["onClick", "icon", "loading"]),
                              createVNode(_component_UButton, {
                                onClick: () => performContainerAction(container, "rebuild"),
                                icon: container.isLoading("rebuild") ? "i-heroicons-spinner" : "i-heroicons-arrow-path-rounded-square",
                                size: "sm",
                                color: "white",
                                square: "",
                                variant: "solid",
                                loading: container.isLoading("rebuild"),
                                label: "Rebuild"
                              }, null, 8, ["onClick", "icon", "loading"]),
                              createVNode(_component_UButton, {
                                to: `/containers/${container.getId()}`,
                                icon: "i-heroicons-eye",
                                size: "sm",
                                color: "white",
                                square: "",
                                variant: "solid",
                                label: "Details"
                              }, null, 8, ["to"])
                            ])
                          ]),
                          _: 2
                        }, 1024);
                      }), 128))
                    ];
                  }
                }),
                _: 2
              }, _parent2, _scopeId));
            } else {
              return [
                createVNode(_component_UPageGrid, { class: "py-4" }, {
                  default: withCtx(() => [
                    (openBlock(true), createBlock(Fragment, null, renderList(item.containers, (container) => {
                      return openBlock(), createBlock(_component_UCard, {
                        key: container.getId()
                      }, {
                        header: withCtx(() => [
                          createTextVNode(toDisplayString(container.getName()) + " - ", 1),
                          createVNode(_component_UBadge, {
                            label: container.getStatus(),
                            color: container.isRunning() ? "green" : "red",
                            variant: "soft"
                          }, null, 8, ["label", "color"])
                        ]),
                        default: withCtx(() => [
                          createVNode("div", { class: "flex flex-wrap gap-2" }, [
                            !container.isRunning() ? (openBlock(), createBlock(_component_UButton, {
                              key: 0,
                              onClick: () => performContainerAction(container, "start"),
                              icon: container.isLoading("start") ? "i-heroicons-spinner" : "i-heroicons-play",
                              size: "sm",
                              color: "white",
                              square: "",
                              variant: "solid",
                              loading: container.isLoading("start"),
                              label: "Start"
                            }, null, 8, ["onClick", "icon", "loading"])) : createCommentVNode("", true),
                            container.isRunning() ? (openBlock(), createBlock(_component_UButton, {
                              key: 1,
                              onClick: () => performContainerAction(container, "restart"),
                              icon: container.isLoading("restart") ? "i-heroicons-spinner" : "i-heroicons-arrow-path",
                              size: "sm",
                              color: "white",
                              square: "",
                              variant: "solid",
                              loading: container.isLoading("restart"),
                              label: "Restart"
                            }, null, 8, ["onClick", "icon", "loading"])) : createCommentVNode("", true),
                            container.isRunning() ? (openBlock(), createBlock(_component_UButton, {
                              key: 2,
                              onClick: () => performContainerAction(container, "stop"),
                              icon: container.isLoading("stop") ? "i-heroicons-spinner" : "i-heroicons-stop",
                              size: "sm",
                              color: "white",
                              square: "",
                              variant: "solid",
                              loading: container.isLoading("stop"),
                              label: "Stop"
                            }, null, 8, ["onClick", "icon", "loading"])) : createCommentVNode("", true),
                            createVNode(_component_UButton, {
                              onClick: () => performContainerAction(container, "remove"),
                              icon: container.isLoading("remove") ? "i-heroicons-spinner" : "i-heroicons-trash",
                              size: "sm",
                              color: "white",
                              square: "",
                              variant: "solid",
                              loading: container.isLoading("remove"),
                              label: "Remove"
                            }, null, 8, ["onClick", "icon", "loading"]),
                            createVNode(_component_UButton, {
                              onClick: () => performContainerAction(container, "rebuild"),
                              icon: container.isLoading("rebuild") ? "i-heroicons-spinner" : "i-heroicons-arrow-path-rounded-square",
                              size: "sm",
                              color: "white",
                              square: "",
                              variant: "solid",
                              loading: container.isLoading("rebuild"),
                              label: "Rebuild"
                            }, null, 8, ["onClick", "icon", "loading"]),
                            createVNode(_component_UButton, {
                              to: `/containers/${container.getId()}`,
                              icon: "i-heroicons-eye",
                              size: "sm",
                              color: "white",
                              square: "",
                              variant: "solid",
                              label: "Details"
                            }, null, 8, ["to"])
                          ])
                        ]),
                        _: 2
                      }, 1024);
                    }), 128))
                  ]),
                  _: 2
                }, 1024)
              ];
            }
          }),
          _: 1
        }, _parent));
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
              _push2(` Loading containers... `);
            } else {
              return [
                createTextVNode(" Loading containers... ")
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
  (ssrContext.modules || (ssrContext.modules = /* @__PURE__ */ new Set())).add("pages/other-containers/index.vue");
  return _sfc_setup ? _sfc_setup(props, ctx) : void 0;
};

export { _sfc_main as default };
//# sourceMappingURL=index-lwer3w52.mjs.map
