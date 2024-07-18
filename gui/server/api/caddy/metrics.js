import Container from "~/utils/Container";

function parseMetrics(metricsText) {
  const metrics = [];
  metricsText.forEach((line) => {
    if (line.startsWith("#")) return; // Ignore comments and metadata
    const match = line.match(/(\w+)\{([^}]*)\}\s+([\d\.]+)/);
    if (match) {
      const [, name, labels, value] = match;
      const labelObj = {};
      labels.split(",").forEach((label) => {
        const [key, val] = label.split("=");
        labelObj[key] = val.replace(/"/g, "");
      });
      metrics.push({ name, labels: labelObj, value: parseFloat(value) });
    }
  });
  return metrics;
}

export default defineEventHandler(async (event) => {
  if (event.node.req.method !== "GET") {
    return { status: "error" };
  }

  try {
    const container = await Container.fetchContainerByName("frankenphp_devner");
    if (!container) {
      return { status: "error", message: "Failed to fetch container" };
    }

    const metricsText = await container.cmd("curl localhost:2019/metrics");
    const parsedMetrics = parseMetrics(metricsText);

    return { status: "success", data: parsedMetrics };
  } catch (error) {
    return { status: "error", message: error.message };
  }
});
