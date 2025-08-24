export namespace main {
	
	export class Connection {
	    ConnectedToOrca: boolean;
	    ConnectionMessage: string;
	    OrcaClient?: pkg.OrcaClient;
	
	    static createFrom(source: any = {}) {
	        return new Connection(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ConnectedToOrca = source["ConnectedToOrca"];
	        this.ConnectionMessage = source["ConnectionMessage"];
	        this.OrcaClient = this.convertValues(source["OrcaClient"], pkg.OrcaClient);
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

export namespace pkg {
	
	export class ObjectNavigator {
	
	
	    static createFrom(source: any = {}) {
	        return new ObjectNavigator(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class FlatReviewPresenter {
	
	
	    static createFrom(source: any = {}) {
	        return new FlatReviewPresenter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class TableNavigator {
	
	
	    static createFrom(source: any = {}) {
	        return new TableNavigator(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class ClipboardPresenter {
	
	
	    static createFrom(source: any = {}) {
	        return new ClipboardPresenter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class NotificationPresenter {
	
	
	    static createFrom(source: any = {}) {
	        return new NotificationPresenter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class SystemInformationPresenter {
	
	
	    static createFrom(source: any = {}) {
	        return new SystemInformationPresenter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class SleepModeManager {
	
	
	    static createFrom(source: any = {}) {
	        return new SleepModeManager(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class StructuralNavigator {
	
	
	    static createFrom(source: any = {}) {
	        return new StructuralNavigator(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class WhereAmIPresenter {
	
	
	    static createFrom(source: any = {}) {
	        return new WhereAmIPresenter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class ActionPresenter {
	
	
	    static createFrom(source: any = {}) {
	        return new ActionPresenter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class SpeechAndVerbosityManager {
	
	
	    static createFrom(source: any = {}) {
	        return new SpeechAndVerbosityManager(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class OrcaClient {
	    // Go type: SpeechAndVerbosityManager
	    SpeechAndVerbosityManager: any;
	    // Go type: ActionPresenter
	    ActionPresenter: any;
	    // Go type: WhereAmIPresenter
	    WhereAmIPresenter: any;
	    // Go type: StructuralNavigator
	    StructuralNavigator: any;
	    // Go type: SleepModeManager
	    SleepModeManager: any;
	    // Go type: SystemInformationPresenter
	    SystemInformationPresenter: any;
	    // Go type: NotificationPresenter
	    NotificationPresenter: any;
	    // Go type: ClipboardPresenter
	    ClipboardPresenter: any;
	    // Go type: TableNavigator
	    TableNavigator: any;
	    // Go type: FlatReviewPresenter
	    FlatReviewPresenter: any;
	    // Go type: ObjectNavigator
	    ObjectNavigator: any;
	
	    static createFrom(source: any = {}) {
	        return new OrcaClient(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.SpeechAndVerbosityManager = this.convertValues(source["SpeechAndVerbosityManager"], null);
	        this.ActionPresenter = this.convertValues(source["ActionPresenter"], null);
	        this.WhereAmIPresenter = this.convertValues(source["WhereAmIPresenter"], null);
	        this.StructuralNavigator = this.convertValues(source["StructuralNavigator"], null);
	        this.SleepModeManager = this.convertValues(source["SleepModeManager"], null);
	        this.SystemInformationPresenter = this.convertValues(source["SystemInformationPresenter"], null);
	        this.NotificationPresenter = this.convertValues(source["NotificationPresenter"], null);
	        this.ClipboardPresenter = this.convertValues(source["ClipboardPresenter"], null);
	        this.TableNavigator = this.convertValues(source["TableNavigator"], null);
	        this.FlatReviewPresenter = this.convertValues(source["FlatReviewPresenter"], null);
	        this.ObjectNavigator = this.convertValues(source["ObjectNavigator"], null);
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

