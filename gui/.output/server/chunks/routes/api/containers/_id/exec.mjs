import { d as defineEventHandler, r as readBody } from '../../../../runtime.mjs';
import Docker from 'dockerode';
import { PassThrough } from 'stream';
import 'node:http';
import 'node:https';
import 'events';
import 'https';
import 'http';
import 'net';
import 'tls';
import 'crypto';
import 'url';
import 'zlib';
import 'buffer';
import 'fs';
import 'path';
import 'engine.io';
import 'socket.io';
import 'node:fs';
import 'node:url';

const docker = new Docker({ socketPath: "/var/run/docker.sock" });
const exec = defineEventHandler(async (event) => {
  if (event.node.req.method !== "POST") {
    return { status: "error" };
  }
  const { id } = event.context.params;
  const body = await readBody(event);
  const command = body.command;
  if (!command) {
    return { status: "error", message: "Command is required" };
  }
  try {
    const container = docker.getContainer(id);
    const exec = await container.exec({
      AttachStdout: true,
      AttachStderr: true,
      Tty: false,
      Cmd: ["sh", "-c", command]
    });
    const stream = await exec.start({ hijack: true });
    const stdout = new PassThrough();
    const stderr = new PassThrough();
    container.modem.demuxStream(stream, stdout, stderr);
    let stdoutOutput = "";
    let stderrOutput = "";
    stdout.on("data", (chunk) => {
      stdoutOutput += chunk.toString("utf8");
    });
    stderr.on("data", (chunk) => {
      stderrOutput += chunk.toString("utf8");
    });
    return new Promise((resolve, reject) => {
      stream.on("end", () => {
        resolve({
          stdout: stdoutOutput,
          stderr: stderrOutput
        });
      });
      stream.on("error", (error) => {
        reject({ status: "error", message: error.message });
      });
    });
  } catch (error) {
    return { status: "error", message: error.message };
  }
});

export { exec as default };
//# sourceMappingURL=exec.mjs.map
