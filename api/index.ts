// Vercel only deploys serverless functions contained within the api directory.
// Rather than nest everything here, this is a hack (which also requires all
// imports in the codebase be relative as Vercel's serverless function build
// process does not respect aliases)
export { default } from "../index";
