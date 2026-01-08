import { useEffect } from 'react';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import { useHistory, useLocation } from '@docusaurus/router';

export default function Home() {
  const { siteConfig, i18n } = useDocusaurusContext();
  const history = useHistory();
  const location = useLocation();
  
  useEffect(() => {
    // Get current locale from URL path
    const pathParts = location.pathname.split('/').filter(Boolean);
    const pathLocale = pathParts[0];
    const currentLocale = i18n.locales.includes(pathLocale) ? pathLocale : i18n.defaultLocale;
    const localePrefix = currentLocale === i18n.defaultLocale ? '' : `${currentLocale}/`;
    
    // Redirect to intro docs with proper locale
    // Use baseUrl from config and redirect to docs/intro (without version prefix)
    // Docusaurus will automatically route to the current version
    const baseUrl = siteConfig.baseUrl || '/';
    // Construct path: baseUrl + localePrefix + docs/intro
    // baseUrl is '/acloud-cli/', localePrefix is '' or 'it/', so result is '/acloud-cli/docs/intro' or '/acloud-cli/it/docs/intro'
    const targetPath = `${baseUrl}${localePrefix}docs/intro`;
    history.replace(targetPath);
  }, [siteConfig, i18n, history, location]);
  
  return null;
}

