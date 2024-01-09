import { ModuleConfig } from "/@/v";
import { useDict } from "./index";

export default (): ModuleConfig => {
	return {
		onLoad({ hasToken }) {
			const { dict } = useDict();
			hasToken(() => {
				dict.refresh();
			});
		}
	};
};
