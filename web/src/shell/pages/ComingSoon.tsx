import { useTranslation } from 'react-i18next'
import type { ModuleLabelKey } from '../navigation'

/** Chỗ giữ cho màn hình chưa xây dựng; `screen` là mã màn hình tại CR-DOC-11 mục 6. */
export function ComingSoon({ labelKey, screen }: { labelKey: ModuleLabelKey; screen: string }) {
  const { t } = useTranslation(['shell', 'runtime', 'content', 'learning'])
  return <p className="border border-dashed border-rule bg-surface p-6 text-muted-foreground">{t('pages.comingSoon', { name: t(labelKey), screen })}</p>
}
