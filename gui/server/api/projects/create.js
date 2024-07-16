import Container from "~/utils/Container.js";
import Project from "~/utils/Project.js";

export default defineEventHandler(async (event) => {
    if (event.node.req.method !== 'POST') {
        return { status: 'error' };
    }

    const body = await readBody(event)
    const { name, type } = body;

    const container = await Container.fetchContainerByName("frankenphp_devner");
    switch (type) {
        case "laravel":
            await container.cmd(`composer create-project laravel/laravel /var/www/html/${name}`);
            break;
        case "wordpress":
            await container.cmd(`wp core download --path=/var/www/html/${name} --allow-root`);
            break;
        default:
            break;
    }

    // watch if the project has been created
    const project = await Project.fetchProject(name);
    
    if (!project) {
        return { status: 'error', message: 'Failed to create project' };
    }

    return { status: 'success' };
});
