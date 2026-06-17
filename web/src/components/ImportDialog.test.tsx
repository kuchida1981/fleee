import { render, screen } from '@testing-library/react';
import { ImportDialog } from './ImportDialog';

describe('ImportDialog', () => {
  it('renders dialog content', () => {
    render(<ImportDialog open={true} onOpenChange={vi.fn()} onImportComplete={vi.fn()} />);

    expect(screen.getByText('勘定科目のインポート')).toBeInTheDocument();
    expect(screen.getByText('インポートを実行')).toBeInTheDocument();
    expect(screen.getByText('CSVまたはTSVファイルを選択してください。')).toBeInTheDocument();
  });

  it('disables import button when no file is selected', () => {
    render(<ImportDialog open={true} onOpenChange={vi.fn()} onImportComplete={vi.fn()} />);

    const button = screen.getByText('インポートを実行');
    expect(button).toBeDisabled();
  });
});
