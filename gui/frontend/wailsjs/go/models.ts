export namespace app {
	
	export class CreateProjectRequest {
	    Name: string;
	    Type: string;
	    DBEngine: string;
	    Template: string;
	    WPInstall: boolean;
	    WPTitle: string;
	    WPAdmin: string;
	    WPPassword: string;
	    WPEmail: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateProjectRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Type = source["Type"];
	        this.DBEngine = source["DBEngine"];
	        this.Template = source["Template"];
	        this.WPInstall = source["WPInstall"];
	        this.WPTitle = source["WPTitle"];
	        this.WPAdmin = source["WPAdmin"];
	        this.WPPassword = source["WPPassword"];
	        this.WPEmail = source["WPEmail"];
	    }
	}
	export class CreateProjectResult {
	    Project: store.Project;
	    DBCreds?: database.Credentials;
	
	    static createFrom(source: any = {}) {
	        return new CreateProjectResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Project = this.convertValues(source["Project"], store.Project);
	        this.DBCreds = this.convertValues(source["DBCreds"], database.Credentials);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace database {
	
	export class Credentials {
	    Host: string;
	    Port: number;
	    Database: string;
	    User: string;
	    Password: string;
	
	    static createFrom(source: any = {}) {
	        return new Credentials(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Host = source["Host"];
	        this.Port = source["Port"];
	        this.Database = source["Database"];
	        this.User = source["User"];
	        this.Password = source["Password"];
	    }
	}

}

export namespace main {
	
	export class ProviderDTO {
	    base_url: string;
	    model: string;
	    kind: string;
	    api_key_env: string;
	    has_key: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ProviderDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.base_url = source["base_url"];
	        this.model = source["model"];
	        this.kind = source["kind"];
	        this.api_key_env = source["api_key_env"];
	        this.has_key = source["has_key"];
	    }
	}
	export class ConfigInfo {
	    config_path: string;
	    data_dir: string;
	    projects_dir: string;
	    active_provider: string;
	    providers: Record<string, ProviderDTO>;
	
	    static createFrom(source: any = {}) {
	        return new ConfigInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.config_path = source["config_path"];
	        this.data_dir = source["data_dir"];
	        this.projects_dir = source["projects_dir"];
	        this.active_provider = source["active_provider"];
	        this.providers = this.convertValues(source["providers"], ProviderDTO, true);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ContainerStatus {
	    name: string;
	    image: string;
	    state: string;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new ContainerStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.image = source["image"];
	        this.state = source["state"];
	        this.status = source["status"];
	    }
	}
	export class CreateDatabaseResult {
	    engine: string;
	    database: string;
	    user: string;
	    password: string;
	    host: string;
	    port: number;
	
	    static createFrom(source: any = {}) {
	        return new CreateDatabaseResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.engine = source["engine"];
	        this.database = source["database"];
	        this.user = source["user"];
	        this.password = source["password"];
	        this.host = source["host"];
	        this.port = source["port"];
	    }
	}
	export class DatabaseInfo {
	    engine: string;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new DatabaseInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.engine = source["engine"];
	        this.name = source["name"];
	    }
	}
	export class DevServerStatusDTO {
	    name: string;
	    running: boolean;
	    pid: number;
	    port: number;
	
	    static createFrom(source: any = {}) {
	        return new DevServerStatusDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.running = source["running"];
	        this.pid = source["pid"];
	        this.port = source["port"];
	    }
	}
	export class HostEntryDTO {
	    domain: string;
	    target: string;
	    managed_at: number;
	
	    static createFrom(source: any = {}) {
	        return new HostEntryDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.domain = source["domain"];
	        this.target = source["target"];
	        this.managed_at = source["managed_at"];
	    }
	}

}

export namespace network {
	
	export class CertsStatus {
	    MkcertInstalled: boolean;
	    CARootInstalled: boolean;
	    CAROOTPath: string;
	
	    static createFrom(source: any = {}) {
	        return new CertsStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.MkcertInstalled = source["MkcertInstalled"];
	        this.CARootInstalled = source["CARootInstalled"];
	        this.CAROOTPath = source["CAROOTPath"];
	    }
	}

}

export namespace store {
	
	export class Project {
	    Name: string;
	    Type: string;
	    DBEngine: string;
	    DBName: string;
	    Domain: string;
	    Path: string;
	    // Go type: time
	    CreatedAt: any;
	    DevPort: number;
	    DevMode: string;
	    DevCommand: string;
	
	    static createFrom(source: any = {}) {
	        return new Project(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Type = source["Type"];
	        this.DBEngine = source["DBEngine"];
	        this.DBName = source["DBName"];
	        this.Domain = source["Domain"];
	        this.Path = source["Path"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.DevPort = source["DevPort"];
	        this.DevMode = source["DevMode"];
	        this.DevCommand = source["DevCommand"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

