/// <reference types="vite/client" />

import * as fs from "node:fs";
import * as path from "node:path";

export const start = () => {
  console.log("Routes started");
  update();
};

export const update = async () => {
  const files = readDirectory(path.join("src", "pages"));

  // Separate layout files from regular page files
  const layoutFiles = files.filter((f) => f.name.includes("layout"));
  const pageFiles = files.filter((f) => !f.name.includes("layout"));

  let output_ts = 'import type {RouteObject} from "react-router-dom";\n';
  output_ts += `const routes: RouteObject[] = [];\n`;

  // Group pages by their directory to apply layouts
  const pagesByDir = new Map<string, typeof pageFiles>();
  pageFiles.forEach((file) => {
    const dir = file.name.includes("/") ? file.name.split("/")[0] : "";
    if (!pagesByDir.has(dir)) {
      pagesByDir.set(dir, []);
    }
    pagesByDir.get(dir)!.push(file);
  });

  // Process each directory group
  pagesByDir.forEach((pagesInDir, dir) => {
    const layoutFile = layoutFiles.find((f) => {
      const layoutDir = f.name.includes("/") ? f.name.split("/")[0] : "";
      return layoutDir === dir;
    });

    if (layoutFile && dir) {
      // This directory has a layout - create parent route with children
      const layoutImportName = getImportName(layoutFile.name.split(".")[0]);
      const layoutPath =
        "./src/pages/" +
        layoutFile.name.split(".")[0].replace(new RegExp("\\\\", "g"), "/");

      output_ts += `import ${layoutImportName} from "${layoutPath}";\n`;

      // Build the parent route path
      const parentPath = getRoutePath("/" + dir);

      // Import all child pages
      const childImports: string[] = [];
      pagesInDir.forEach((file) => {
        const no_ext = file.name.split(".")[0];
        const import_name = getImportName(no_ext);
        const filePath =
          "./src/pages/" + no_ext.replace(new RegExp("\\\\", "g"), "/");
        output_ts += `import ${import_name} from "${filePath}";\n`;
        childImports.push(import_name);
      });

      // Create parent route with layout
      output_ts += `routes.push({\n`;
      output_ts += `  path: "${parentPath}",\n`;
      output_ts += `  element: <${layoutImportName} />,\n`;
      output_ts += `  children: [\n`;

      // Add child routes
      pagesInDir.forEach((file) => {
        const no_ext = file.name.split(".")[0];
        const import_name = getImportName(no_ext);

        // Get the relative path within the directory
        const relativePath = no_ext.replace(dir + "/", "");
        let childPath = getRoutePath("/" + relativePath);

        // Handle index routes
        if (relativePath === "index" || relativePath.endsWith("/index")) {
          output_ts += `    { index: true, element: <${import_name} /> },\n`;
        } else {
          // Remove leading slash for child routes
          childPath = childPath.startsWith("/")
            ? childPath.slice(1)
            : childPath;
          output_ts += `    { path: "${childPath}", element: <${import_name} /> },\n`;
        }
      });

      output_ts += `  ],\n`;
      output_ts += `});\n`;
    } else {
      // No layout - process pages individually (root level pages)
      pagesInDir.forEach((file) => {
        const no_ext = file.name.split(".")[0];
        const import_name = getImportName(no_ext);
        const route_name = getRoutePath("/" + no_ext);
        const filePath =
          "./src/pages/" + no_ext.replace(new RegExp("\\\\", "g"), "/");

        output_ts += `import ${import_name} from "${filePath}";\n`;
        output_ts += `routes.push({ path: "${route_name}", element: <${import_name} /> });\n`;
      });
    }
  });

  output_ts += `export default routes;\n`;
  await fs.writeFileSync(path.join("routes.tsx"), output_ts);
};

const getImportName = (pathStr: string): string => {
  return pathStr
    .replace(new RegExp("\\.", "g"), "_")
    .replace(new RegExp("/", "g"), "_")
    .replace(new RegExp("-", "g"), "")
    .replace(new RegExp("\\[", "g"), "")
    .replace(new RegExp("]", "g"), "")
    .split("_")
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join("");
};

const getRoutePath = (pathStr: string): string => {
  const route_name = pathStr
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

  return route_name_parts.join("/");
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
