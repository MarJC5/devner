import { d as defineEventHandler } from '../../../../runtime.mjs';
import { d as dockerService } from '../../../../_/docker.mjs';
import 'node:http';
import 'node:https';
import 'fs';
import 'path';
import 'node:fs';
import 'node:url';
import 'dockerode';

const remove = defineEventHandler(async (event) => {
  if (event.request.method !== "POST") {
    return { status: "error" };
  }
  const { id } = event.context.params;
  await dockerService.removeContainer(id);
  return { status: "removed" };
});

export { remove as default };
//# sourceMappingURL=remove.mjs.map
