import { Plugin } from "vite";
import { parseJson } from "./utils";
import { createEps, createMenu, createSvg, createTag, getEps, getModules } from "./lib";

export function v(): Plugin {
	return {
		name: "vite-v",
		enforce: "pre",
		configureServer(server) {
			server.middlewares.use(async (req, res, next) => {
				function done(data: any) {
					res.writeHead(200, { "Content-Type": "text/html;charset=UTF-8" });
					res.end(JSON.stringify(data));
				}

				if (req.url?.includes("__v")) {
					const body = await parseJson(req);
					let next: any;

					switch (req.url) {
						// 快速创建菜单
						case "/__v_createMenu":
							next = createMenu(body);
							break;

						// 获取模块列表
						case "/__v_modules":
							next = getModules();
							break;

						// 创建描述文件
						case "/__v_eps":
							next = createEps(body);
							break;
					}

					if (next) {
						next.then((data: any) => {
							done({
								code: 1000,
								data
							});
						}).catch((err: Error) => {
							done({
								code: 1001,
								message: err.message
							});
						});
					}
				} else {
					next();
				}
			});
		},
		transform(code, id) {
			return createTag(code, id);
		},
		transformIndexHtml(html) {
			return createSvg(html);
		},
		config() {
			return {
				define: {
					__EPS__: getEps()
				}
			};
		}
	};
}
