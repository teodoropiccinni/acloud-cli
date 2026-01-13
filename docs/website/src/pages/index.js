import {Redirect} from '@docusaurus/router';

export default function Home() {
  // Redirect root to intro page
  // Docusaurus will handle version routing automatically
  // With routeBasePath: '/', docs are at root, so redirect to '/intro'
  // Docusaurus will automatically:
  // - Add baseUrl (/acloud-cli/)
  // - Add version prefix (/next/) if versioning is enabled
  // - Handle locale routing if needed
  return <Redirect to="/intro" />;
}

