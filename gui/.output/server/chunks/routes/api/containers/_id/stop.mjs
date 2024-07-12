import { d as defineEventHandler } from '../../../../runtime.mjs';
import { d as dockerService } from '../../../../_/docker.mjs';
import 'node:http';
import 'node:https';
import 'fs';
import 'path';
import 'node:fs';
import 'node:url';
import 'dockerode';

const stop = defineEventHandler(async (event) => {
  if (event.request.method !== "POST") {
    return { status: "error" };
  }
  const { id } = event.context.params;
  await dockerService.stopContainer(id);
  return { status: "stopped" };
});

export { stop as default };
//# sourceMappingURL=stop.mjs.map
