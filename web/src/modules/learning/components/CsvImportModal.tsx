import { useState, type ChangeEvent } from 'react'
import { useTranslation } from 'react-i18next'
import { Upload, AlertTriangle, CheckCircle, X } from 'lucide-react'
import { Button } from '@/shell/ui/button'
import { Card, CardTitle } from '@/shell/ui/card'
import { importCsv, type CsvPreviewResponse } from '../api/classApi'

interface Props {
  classId: string
  isOpen: boolean
  onClose: () => void
  onSuccess: () => void
}

export function CsvImportModal({ classId, isOpen, onClose, onSuccess }: Props) {
  const { t } = useTranslation()
  const [file, setFile] = useState<File | null>(null)
  const [preview, setPreview] = useState<CsvPreviewResponse | null>(null)
  const [loading, setLoading] = useState(false)
  const [errorMsg, setErrorMsg] = useState<string | null>(null)

  if (!isOpen) return null

  const handleFileChange = (e: ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      setFile(e.target.files[0])
      setPreview(null)
      setErrorMsg(null)
    }
  }

  const handlePreview = async () => {
    if (!file) return
    setLoading(true)
    setErrorMsg(null)
    try {
      const res = await importCsv(classId, file, false)
      setPreview(res as CsvPreviewResponse)
    } catch (err: unknown) {
      const msg = typeof err === 'object' && err && 'detail' in err ? String((err as { detail: string }).detail) : 'Lỗi phân tích tệp CSV'
      setErrorMsg(msg)
    } finally {
      setLoading(false)
    }
  }

  const handleCommit = async () => {
    if (!file) return
    setLoading(true)
    setErrorMsg(null)
    try {
      await importCsv(classId, file, true)
      onSuccess()
      onClose()
    } catch (err: unknown) {
      const msg = typeof err === 'object' && err && 'detail' in err ? String((err as { detail: string }).detail) : 'Lỗi lưu dữ liệu sinh viên'
      setErrorMsg(msg)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4">
      <Card className="w-full max-w-2xl bg-surface border-rule max-h-[90vh] overflow-y-auto">
        <div className="flex items-center justify-between border-b border-rule pb-3">
          <CardTitle>{t('learning:csv.title')}</CardTitle>
          <button onClick={onClose} className="text-muted-foreground hover:text-foreground">
            <X className="size-5" />
          </button>
        </div>

        <div className="mt-4 space-y-4">
          <label className="flex flex-col items-center justify-center border-2 border-dashed border-rule p-6 hover:bg-accent/40 cursor-pointer">
            <Upload className="size-8 text-primary mb-2" />
            <span className="text-sm font-medium">{file ? file.name : t('learning:csv.dragDrop')}</span>
            <input type="file" accept=".csv" onChange={handleFileChange} className="hidden" />
          </label>

          {errorMsg && (
            <div className="border border-danger bg-danger-soft p-3 text-sm text-danger flex items-center gap-2">
              <AlertTriangle className="size-4 shrink-0" />
              <span>{errorMsg}</span>
            </div>
          )}

          {file && !preview && (
            <div className="flex justify-end gap-2">
              <Button variant="outline" onClick={onClose}>{t('learning:classes.cancel')}</Button>
              <Button variant="primary" disabled={loading} onClick={handlePreview}>
                {t('learning:csv.preview')}
              </Button>
            </div>
          )}

          {preview && (
            <div className="space-y-4">
              <div className="flex items-center gap-4 text-sm font-mono">
                <span className="flex items-center gap-1.5 text-success">
                  <CheckCircle className="size-4" />
                  {t('learning:csv.valid')}: {preview.validRows} / {preview.totalRows}
                </span>
                {preview.errorRows > 0 && (
                  <span className="flex items-center gap-1.5 text-danger font-bold">
                    <AlertTriangle className="size-4" />
                    {t('learning:csv.errors')}: {preview.errorRows}
                  </span>
                )}
              </div>

              {preview.errors.length > 0 && (
                <div className="border border-signal bg-signal-soft/20 p-3 max-h-48 overflow-y-auto">
                  <table className="w-full text-xs font-mono">
                    <thead>
                      <tr className="border-b border-signal/40 text-left">
                        <th className="pb-1">{t('learning:csv.line')}</th>
                        <th className="pb-1">{t('learning:csv.field')}</th>
                        <th className="pb-1">{t('learning:csv.message')}</th>
                      </tr>
                    </thead>
                    <tbody>
                      {preview.errors.map((err, idx) => (
                        <tr key={idx} className="border-b border-signal/20">
                          <td className="py-1">{err.line}</td>
                          <td className="py-1">{err.field}</td>
                          <td className="py-1">{err.message}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              )}

              <div className="flex justify-end gap-2 pt-2">
                <Button variant="outline" onClick={onClose}>{t('learning:classes.cancel')}</Button>
                <Button
                  variant="primary"
                  disabled={loading || preview.errors.length > 0}
                  onClick={handleCommit}
                >
                  {t('learning:csv.commit')}
                </Button>
              </div>
            </div>
          )}
        </div>
      </Card>
    </div>
  )
}
