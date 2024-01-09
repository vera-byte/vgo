import { ModuleConfig } from "/@/v";

export default (): ModuleConfig => {
	return {
		components: [import("./components/export-btn")]
	};
};
