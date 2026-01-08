import { useEffect } from 'react';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import { useHistory, useLocation } from '@docusaurus/router';

export default function Home() {
  const { i18n } = useDocusaurusContext();
  const history = useHistory();
  const location = useLocation();
  
  useEffect(() => {
    // Get current locale from URL path
    const pathParts = location.pathname.split('/').filter(Boolean);
    const pathLocale = pathParts[0];
    const currentLocale = i18n.locales.includes(pathLocale) ? pathLocale : i18n.defaultLocale;
    const localePrefix = currentLocale === i18n.defaultLocale ? '' : `${currentLocale}/`;
    
    // Redirect to intro docs with proper locale
    // Docusaurus router automatically handles baseUrl
    // Version path "next" is the version path configured
    const targetPath = `/${localePrefix}docs/next/intro`;
    history.replace(targetPath);
  }, [i18n, history, location]);
  
  return null;
}

