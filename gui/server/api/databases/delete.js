import Container from "@/models/Container";
import Database from "@/models/Database";

export default defineEventHandler(async (event) => {
    if (event.node.req.method !== 'POST') {
        return { status: 'error' };
    }

    const body = await readBody(event)
    const { type, username, password, database } = body;

    const container = await Container.fetchContainerByName(`${type}_devner`);
    await container.cmd(`mysql -u ${username} -p'${password}' -e "DROP DATABASE IF EXISTS ${database}"`);

    // watch if the project has been deleted
    const db = await Database.fetchDatabase(database);
    
    if (db) {
        return { status: 'error', message: 'Failed to delete database' };
    }

    return { status: 'success' };
});
