import { useState, type ChangeEvent } from 'react'
import { useTranslation } from 'react-i18next'
import { Upload, AlertTriangle, CheckCircle, X } from 'lucide-react'
import type { ApiProblem } from '@/shell/i18n/problem'
import { Button } from '@/shell/ui/button'
import { Card, CardTitle } from '@/shell/ui/card'
import { ProblemAlert } from '@/shell/ui/problem-alert'
import { previewCsvMembers, importCsvMembers, type CsvPreviewResponse } from '../api/classApi'

interface Props {
  classId: string
  isOpen: boolean
  onClose: () => void
  onSuccess: (importedCount: number) => void
}

export function CsvImportModal({ classId, isOpen, onClose, onSuccess }: Props) {
  const { t } = useTranslation(['learning', 'shell'])
  const [file, setFile] = useState<File | null>(null)
  const [preview, setPreview] = useState<CsvPreviewResponse | null>(null)
  const [loading, setLoading] = useState(false)
  const [errorProblem, setErrorProblem] = useState<ApiProblem | null>(null)

  if (!isOpen) return null

  const handleFileChange = (e: ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      setFile(e.target.files[0])
      setPreview(null)
      setErrorProblem(null)
    }
  }

  const handlePreview = async () => {
    if (!file) return
    setLoading(true)
    setErrorProblem(null)
    try {
      const res = await previewCsvMembers(classId, file)
      setPreview(res)
    } catch (err: unknown) {
      setErrorProblem(typeof err === 'object' && err ? (err as ApiProblem) : {})
    } finally {
      setLoading(false)
    }
  }

  const handleCommit = async () => {
    if (!file) return
    setLoading(true)
    setErrorProblem(null)
    try {
      const res = await importCsvMembers(classId, file)
      onSuccess(res.importedCount)
      onClose()
    } catch (err: unknown) {
      setErrorProblem(typeof err === 'object' && err ? (err as ApiProblem) : {})
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

          {errorProblem && <ProblemAlert problem={errorProblem} ns="learning" />}

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
