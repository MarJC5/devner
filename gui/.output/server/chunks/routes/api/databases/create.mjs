import { d as defineEventHandler, r as readBody } from '../../../runtime.mjs';
import { C as Container } from '../../../_/Container.mjs';
import { D as Database } from '../../../_/Database.mjs';
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
import 'strip-ansi';

const create = defineEventHandler(async (event) => {
  if (event.node.req.method !== "POST") {
    return { status: "error" };
  }
  const body = await readBody(event);
  const { type, username, password, database } = body;
  const container = await Container.fetchContainerByName(`${type}_devner`);
  await container.cmd(`mysql -u ${username} -p'${password}' -e "CREATE DATABASE IF NOT EXISTS ${database}"`);
  const db = await Database.fetchDatabase(database);
  if (!db) {
    return { status: "error", message: "Failed to create database" };
  }
  return { status: "success" };
});

export { create as default };
//# sourceMappingURL=create.mjs.map
