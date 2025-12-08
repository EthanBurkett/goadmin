/// <reference types="vite/client" />

import * as fs from "node:fs";
import * as path from "node:path";

export const start = () => {
  console.log("Routes started");
  update();
};

export const update = async () => {
  const files = readDirectory(path.join("src", "pages"));

  let output_ts = 'import {RouteObject} from "react-router-dom";\n';
  output_ts += `const routes: RouteObject[] = [];\n`;
  files.forEach(async (file) => {
    const no_ext = file.name.split(".")[0];
    let import_name = no_ext
      .replace(new RegExp("\\.", "g"), "_")
      .replace(new RegExp("/", "g"), "_")
      .replace(new RegExp("-", "g"), "");
    let route_name =
      "/" +
      no_ext
        .replace(new RegExp("\\.", "g"), "_")
        .replace(new RegExp("/", "g"), "_")
        .split("_")
        .join("/")
        .replace(/index/g, "");

    const route_name_parts = route_name.split("/");
    route_name_parts.forEach((part, index) => {
      if (part.startsWith("[") && part.endsWith("]")) {
        route_name_parts[index] = ":" + part.slice(1, part.length - 1);
      }
    });
    route_name = route_name_parts.join("/");

    import_name = import_name
      .replace(new RegExp("\\[", "g"), "")
      .replace(new RegExp("]", "g"), "")
      .split("_")
      .map((word) => {
        return word.charAt(0).toUpperCase() + word.slice(1);
      })
      .join("");

    const path = "./src/pages/" + no_ext.replace(new RegExp("\\\\", "g"), "/");

    output_ts += `import ${import_name} from "${path}";\n`;
    output_ts += `routes.push({ path: "${route_name}", element: <${import_name} /> });\n`;
  });

  output_ts += `export default routes;\n`;
  await fs.writeFileSync(path.join("routes.tsx"), output_ts);
};

const readDirectory = (dir: string) => {
  const files: {
    name: string;
    path: string;
  }[] = [];
  fs.readdirSync(dir, {
    withFileTypes: true,
  }).forEach((file) => {
    if (file.isDirectory()) {
      readDirectory(path.join(dir, file.name)).forEach((subfile) => {
        files.push({
          name: file.name + "/" + subfile.name,
          path: file.name + "/" + subfile.path,
        });
      });
    } else {
      files.push({
        name: file.name,
        path: file.name,
      });
    }
  });

  return files;
};
