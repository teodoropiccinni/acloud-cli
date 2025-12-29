import { useEffect } from 'react';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';

export default function Home() {
  const { siteConfig } = useDocusaurusContext();
  
  useEffect(() => {
    const baseUrl = siteConfig.baseUrl || '/';
    // Ensure baseUrl ends with / and construct target path
    const normalizedBaseUrl = baseUrl.endsWith('/') ? baseUrl : `${baseUrl}/`;
    const targetPath = `${normalizedBaseUrl}docs/intro`;
    // Always redirect from homepage to intro docs
    window.location.replace(targetPath);
  }, [siteConfig]);
  
  return null;
}

