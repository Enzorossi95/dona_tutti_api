-- +goose Up
-- Campaign Closure Reports table
CREATE TABLE IF NOT EXISTS campaign_closure_reports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    campaign_id UUID NOT NULL REFERENCES campaigns(id) ON DELETE CASCADE,

    -- Closure type
    closure_type VARCHAR(50) NOT NULL CHECK (closure_type IN ('goal_reached', 'end_date', 'manual')),
    closure_reason TEXT,
    closed_by UUID REFERENCES users(id),

    -- Financial metrics
    total_raised DECIMAL(12,2) NOT NULL DEFAULT 0,
    total_donors INTEGER NOT NULL DEFAULT 0,
    total_donations INTEGER NOT NULL DEFAULT 0,
    campaign_goal DECIMAL(12,2) NOT NULL,
    goal_percentage DECIMAL(5,2) NOT NULL DEFAULT 0,

    -- Expenses
    total_expenses DECIMAL(12,2) NOT NULL DEFAULT 0,
    total_receipts INTEGER NOT NULL DEFAULT 0,
    receipts_with_documents INTEGER NOT NULL DEFAULT 0,

    -- Activities
    total_activities INTEGER NOT NULL DEFAULT 0,

    -- Transparency
    transparency_score DECIMAL(5,2) NOT NULL DEFAULT 0 CHECK (transparency_score >= 0 AND transparency_score <= 100),
    transparency_breakdown JSONB,

    -- Alerts (placeholder for future implementation)
    alerts_count INTEGER NOT NULL DEFAULT 0,
    alerts_resolved INTEGER NOT NULL DEFAULT 0,

    -- PDF report
    report_pdf_url TEXT,
    report_hash VARCHAR(64),

    -- Timestamps
    closed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- One report per campaign
    CONSTRAINT unique_campaign_closure UNIQUE(campaign_id)
);

-- Indexes for campaign_closure_reports
CREATE INDEX IF NOT EXISTS idx_closure_reports_campaign_id ON campaign_closure_reports(campaign_id);
CREATE INDEX IF NOT EXISTS idx_closure_reports_closed_at ON campaign_closure_reports(closed_at);
CREATE INDEX IF NOT EXISTS idx_closure_reports_transparency_score ON campaign_closure_reports(transparency_score);

-- Campaign Alerts table (placeholder for future implementation)
CREATE TABLE IF NOT EXISTS campaign_alerts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    campaign_id UUID NOT NULL REFERENCES campaigns(id) ON DELETE CASCADE,
    alert_type VARCHAR(50) NOT NULL,
    description TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'investigating', 'resolved', 'dismissed')),
    severity VARCHAR(20) NOT NULL DEFAULT 'medium' CHECK (severity IN ('low', 'medium', 'high', 'critical')),
    reported_by UUID REFERENCES users(id),
    resolved_by UUID REFERENCES users(id),
    resolution_notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    resolved_at TIMESTAMP WITH TIME ZONE
);

-- Indexes for campaign_alerts
CREATE INDEX IF NOT EXISTS idx_campaign_alerts_campaign_id ON campaign_alerts(campaign_id);
CREATE INDEX IF NOT EXISTS idx_campaign_alerts_status ON campaign_alerts(status);
CREATE INDEX IF NOT EXISTS idx_campaign_alerts_severity ON campaign_alerts(severity);

-- Comments
COMMENT ON TABLE campaign_closure_reports IS 'Stores closure/audit reports for campaigns';
COMMENT ON TABLE campaign_alerts IS 'Placeholder table for future alert system implementation';
COMMENT ON COLUMN campaign_closure_reports.transparency_score IS 'Score from 0-100 based on transparency metrics';
COMMENT ON COLUMN campaign_closure_reports.transparency_breakdown IS 'JSON with detailed score calculation breakdown';

-- +goose Down
DROP TABLE IF EXISTS campaign_alerts;
DROP TABLE IF EXISTS campaign_closure_reports;
