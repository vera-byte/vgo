import { Merge, ModuleConfig } from "/@/v";
import Crud from "@v-vue/crud";
import "@v-vue/crud/dist/index.css";

export default (): Merge<ModuleConfig, CrudOptions> => {
	return {
		options: {
			dict: {
				sort: {
					prop: "order",
					order: "sort"
				}
			}
		},
		install: Crud.install
	};
};
