import { u as useUI } from './Button-CDYqYkPm.mjs';
import { defineComponent, toRef, mergeProps, unref, useSSRContext, computed, createVNode, resolveDynamicComponent, withCtx, openBlock, createBlock, renderSlot, createCommentVNode } from 'vue';
import { ssrRenderAttrs, ssrRenderSlot, ssrRenderVNode, ssrRenderClass } from 'vue/server-renderer';
import { twMerge, twJoin } from 'tailwind-merge';
import { m as mergeConfig, b as appConfig } from './server.mjs';
import { _ as _export_sfc } from './_plugin-vue_export-helper-1tPrXgE0.mjs';
import stripAnsi from 'strip-ansi';

const card = {
  base: "",
  background: "bg-white dark:bg-gray-900",
  divide: "divide-y divide-gray-200 dark:divide-gray-800",
  ring: "ring-1 ring-gray-200 dark:ring-gray-800",
  rounded: "rounded-lg",
  shadow: "shadow",
  body: {
    base: "",
    background: "",
    padding: "px-4 py-5 sm:p-6"
  },
  header: {
    base: "",
    background: "",
    padding: "px-4 py-5 sm:px-6"
  },
  footer: {
    base: "",
    background: "",
    padding: "px-4 py-4 sm:px-6"
  }
};
const _sfc_main$1 = /* @__PURE__ */ defineComponent({
  ...{
    inheritAttrs: false
  },
  __name: "PageGrid",
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
      wrapper: "grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 gap-8"
    };
    const props = __props;
    const { ui, attrs } = useUI("page.grid", toRef(props, "ui"), config2, toRef(props, "class"), true);
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
  (ssrContext.modules || (ssrContext.modules = /* @__PURE__ */ new Set())).add("node_modules/@nuxt/ui-pro/components/page/PageGrid.vue");
  return _sfc_setup$1 ? _sfc_setup$1(props, ctx) : void 0;
};
const config = mergeConfig(appConfig.ui.strategy, appConfig.ui.card, card);
const _sfc_main = defineComponent({
  inheritAttrs: false,
  props: {
    as: {
      type: String,
      default: "div"
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
    const { ui, attrs } = useUI("card", toRef(props, "ui"), config);
    const cardClass = computed(() => {
      return twMerge(twJoin(
        ui.value.base,
        ui.value.rounded,
        ui.value.divide,
        ui.value.ring,
        ui.value.shadow,
        ui.value.background
      ), props.class);
    });
    return {
      // eslint-disable-next-line vue/no-dupe-keys
      ui,
      attrs,
      cardClass
    };
  }
});
function _sfc_ssrRender(_ctx, _push, _parent, _attrs, $props, $setup, $data, $options) {
  ssrRenderVNode(_push, createVNode(resolveDynamicComponent(_ctx.$attrs.onSubmit ? "form" : _ctx.as), mergeProps({ class: _ctx.cardClass }, _ctx.attrs, _attrs), {
    default: withCtx((_, _push2, _parent2, _scopeId) => {
      if (_push2) {
        if (_ctx.$slots.header) {
          _push2(`<div class="${ssrRenderClass([_ctx.ui.header.base, _ctx.ui.header.padding, _ctx.ui.header.background])}"${_scopeId}>`);
          ssrRenderSlot(_ctx.$slots, "header", {}, null, _push2, _parent2, _scopeId);
          _push2(`</div>`);
        } else {
          _push2(`<!---->`);
        }
        if (_ctx.$slots.default) {
          _push2(`<div class="${ssrRenderClass([_ctx.ui.body.base, _ctx.ui.body.padding, _ctx.ui.body.background])}"${_scopeId}>`);
          ssrRenderSlot(_ctx.$slots, "default", {}, null, _push2, _parent2, _scopeId);
          _push2(`</div>`);
        } else {
          _push2(`<!---->`);
        }
        if (_ctx.$slots.footer) {
          _push2(`<div class="${ssrRenderClass([_ctx.ui.footer.base, _ctx.ui.footer.padding, _ctx.ui.footer.background])}"${_scopeId}>`);
          ssrRenderSlot(_ctx.$slots, "footer", {}, null, _push2, _parent2, _scopeId);
          _push2(`</div>`);
        } else {
          _push2(`<!---->`);
        }
      } else {
        return [
          _ctx.$slots.header ? (openBlock(), createBlock("div", {
            key: 0,
            class: [_ctx.ui.header.base, _ctx.ui.header.padding, _ctx.ui.header.background]
          }, [
            renderSlot(_ctx.$slots, "header")
          ], 2)) : createCommentVNode("", true),
          _ctx.$slots.default ? (openBlock(), createBlock("div", {
            key: 1,
            class: [_ctx.ui.body.base, _ctx.ui.body.padding, _ctx.ui.body.background]
          }, [
            renderSlot(_ctx.$slots, "default")
          ], 2)) : createCommentVNode("", true),
          _ctx.$slots.footer ? (openBlock(), createBlock("div", {
            key: 2,
            class: [_ctx.ui.footer.base, _ctx.ui.footer.padding, _ctx.ui.footer.background]
          }, [
            renderSlot(_ctx.$slots, "footer")
          ], 2)) : createCommentVNode("", true)
        ];
      }
    }),
    _: 3
  }), _parent);
}
const _sfc_setup = _sfc_main.setup;
_sfc_main.setup = (props, ctx) => {
  const ssrContext = useSSRContext();
  (ssrContext.modules || (ssrContext.modules = /* @__PURE__ */ new Set())).add("node_modules/@nuxt/ui/dist/runtime/components/layout/Card.vue");
  return _sfc_setup ? _sfc_setup(props, ctx) : void 0;
};
const __nuxt_component_3 = /* @__PURE__ */ _export_sfc(_sfc_main, [["ssrRender", _sfc_ssrRender]]);
class Container {
  constructor(containerData) {
    this.id = containerData.Id || 0;
    this.created = containerData.Created || "";
    this.name = containerData.Name || "";
    this.status = containerData.State.Status || "";
    this.pid = containerData.State.Pid || 0;
    this.startedAt = containerData.State.StartedAt || "";
    this.image = containerData.Image || "";
    this.ports = containerData.NetworkSettings.Ports || {};
    this.labels = containerData.Config.Labels || {};
    this.network = containerData.NetworkSettings.Networks || {};
    this.project = this.labels["com.docker.compose.project"] || "";
    this.loading = {
      start: false,
      restart: false,
      stop: false,
      remove: false,
      rebuild: false
    };
  }
  static async logs(id, options = {}) {
    const logStream = await $fetch(`/api/containers/${id}/logs`, {
      method: "POST",
      body: JSON.stringify(options),
      headers: { "Content-Type": "application/json" }
    });
    logStream.on("data", (chunk) => {
      return chunk.toString("utf8");
    });
    logStream.on("end", () => {
      console.log("Stream ended");
    });
    logStream.on("error", (error) => {
      console.error("Stream error:", error);
    });
  }
  static async fetchContainer(id) {
    const containerData = await $fetch(`/api/containers/${id}/details`, {
      method: "GET"
    });
    return new Container(containerData);
  }
  static async fetchContainerByName(containerName) {
    try {
      const containers = await Container.all();
      return containers.find(
        (container) => container.name.includes(`/${containerName}`)
      );
    } catch (error) {
      console.error("Error fetching container by name:", error);
      throw error;
    }
  }
  static async all() {
    const data = await $fetch(`/api/containers`, { method: "GET" });
    const containers = await Promise.all(
      data.map(async (containerData) => {
        const model = await Container.fetchContainer(containerData.Id);
        if (model.getNetworks().devner && !model.getName().toLowerCase().includes("gui")) {
          return model;
        }
      })
    ).then((results) => results.filter((container) => container !== void 0));
    return containers;
  }
  static async allOthers() {
    const data = await $fetch(`/api/containers`, { method: "GET" });
    const containers = await Promise.all(
      data.map(async (containerData) => {
        const model = await Container.fetchContainer(containerData.Id);
        if (!model.getNetworks().devner) {
          return model;
        }
      })
    ).then((results) => results.filter((container) => container !== void 0));
    const groupedContainers = containers.reduce((acc, container) => {
      const project = container.project;
      if (!acc[project]) {
        acc[project] = [];
      }
      acc[project].push(container);
      return acc;
    }, {});
    const formattedContainers = Object.entries(groupedContainers).map(([label, containers2]) => ({
      label: label.charAt(0).toUpperCase() + label.slice(1),
      icon: "i-heroicons-cube-transparent",
      containers: containers2
    }));
    return formattedContainers;
  }
  formatDate(isoString) {
    const date = new Date(isoString);
    return date.toLocaleString();
  }
  /**
   * Perform an action on the container
   *
   * @param {string} action
   */
  async performAction(action) {
    this.loading[action] = true;
    console.log(`Performing action ${action} on container with ID: ${this.id}`);
    await $fetch(`/api/containers/${this.id}/${action}`, { method: "POST" });
    this.loading[action] = false;
    await this.updateStatus();
  }
  async start() {
    await this.performAction("start");
  }
  async restart() {
    await this.performAction("restart");
  }
  async stop() {
    await this.performAction("stop");
  }
  async remove() {
    await this.performAction("remove");
  }
  async rebuild() {
    await this.performAction("rebuild");
  }
  async updateStatus() {
    const updatedContainer = await Container.fetchContainer(this.id);
    this.status = updatedContainer.status;
    this.pid = updatedContainer.pid;
    this.startedAt = updatedContainer.startedAt;
  }
  /**
   * Exec a command in the container
   */
  async cmd(command, parseAsTable = false, clean = true) {
    try {
      const response = await $fetch(`/api/containers/${this.id}/exec`, {
        method: "POST",
        body: JSON.stringify({ command }),
        headers: { "Content-Type": "application/json" }
      });
      if (response.status === "error") {
        throw new Error(response.message);
      }
      const rawOutput = response.stdout.split("\n");
      if (parseAsTable) {
        return this.parseTableOutput(rawOutput);
      }
      if (clean) {
        return rawOutput.map((line) => stripAnsi(line).replace(/[^\x20-\x7E]/g, "")).filter((line) => line);
      }
      return rawOutput;
    } catch (error) {
      console.error("Error executing command in container:", error);
      throw error;
    }
  }
  parseTableOutput(output) {
    const [header, ...rows] = output;
    const headers = header.split("	").map((h) => h.trim());
    return rows.map((row) => {
      const columns = row.split("	").map((c) => c.trim());
      let rowObject = {};
      headers.forEach((header2, index) => {
        rowObject[stripAnsi(header2).replace(/[^\x20-\x7E]/g, "")] = columns[index] || "";
      });
      return rowObject;
    });
  }
  /**
   * Get Directories
   */
  async getDirectories(path = "/") {
    const filesAndDirs = await this.cmd(`ls -a ${path}`);
    return { filesAndDirs, path, count: filesAndDirs.length };
  }
  /**
   * Get id of the container
   *
   * @returns {string}
   */
  getId() {
    return this.id;
  }
  /**
   * Get creation date of the container
   *
   * @returns {string}
   */
  getCreated() {
    return this.formatDate(this.created);
  }
  /**
   * Get name of the container
   *
   * @returns {string}
   */
  getName() {
    let friendlyName = this.name.replace("/", "").replace(/[-_]/g, " ").replace(/\b\w/g, (c) => c.toUpperCase());
    friendlyName = friendlyName.replace(/ Devner$/, "");
    return friendlyName;
  }
  /**
   * Get status of the container
   *
   * @returns {string}
   */
  getStatus(format = "capitalize") {
    if (format === "capitalize") {
      return this.status.charAt(0).toUpperCase() + this.status.slice(1);
    }
    return this.status;
  }
  /**
   * Get image of the container
   *
   * @returns {string}
   */
  getImage() {
    return this.image;
  }
  /**
   * Get ports of the container
   *
   * @returns {object}
   */
  getPorts(full = true) {
    const portsArray = [];
    for (const [port, mappings] of Object.entries(this.ports)) {
      if (mappings) {
        mappings.forEach((mapping) => {
          if (full) {
            portsArray.push(`${mapping.HostPort}/${port.split("/")[1]}`);
          } else {
            portsArray.push(mapping.HostPort);
          }
        });
      }
    }
    return portsArray;
  }
  /**
   * Get labels of the container
   *
   * @returns {object}
   */
  getLabels() {
    return this.labels;
  }
  /**
   * Get NetworkSettings of the container
   *
   * @returns {object}
   */
  getNetworks() {
    return this.network;
  }
  /**
   * Is running
   */
  isRunning() {
    return this.status === "running";
  }
  /**
   * Actions status
   */
  isLoading(action) {
    return this.loading[action];
  }
}

export { Container as C, _sfc_main$1 as _, __nuxt_component_3 as a };
//# sourceMappingURL=Container-B3wCur7c.mjs.map
