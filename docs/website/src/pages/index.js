import {Redirect} from '@docusaurus/router';
import useBaseUrl from '@docusaurus/useBaseUrl';

export default function Home() {
  // Redirect root to intro page
  // Use useBaseUrl to ensure baseUrl is properly included in the redirect path
  // Docusaurus will automatically add version prefix (/next/) if versioning is enabled
  const introUrl = useBaseUrl('/intro');
  return <Redirect to={introUrl} />;
}

