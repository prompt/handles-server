// @ts-check

import eslint from "@eslint/js";
import tseslint from "typescript-eslint";

export default tseslint.config(
  { ignores: ["dist/*", ".vercel/*", "bin/*", "builds/*"] },
  eslint.configs.recommended,
  ...tseslint.configs.recommended,
);
