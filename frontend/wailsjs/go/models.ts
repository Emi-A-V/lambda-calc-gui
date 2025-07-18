export namespace main {
	
	export class CalculationResult {
	    Equation: string;
	    Error: string;
	    ErrorID: number;
	
	    static createFrom(source: any = {}) {
	        return new CalculationResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Equation = source["Equation"];
	        this.Error = source["Error"];
	        this.ErrorID = source["ErrorID"];
	    }
	}

}

