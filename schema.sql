--
-- PostgreSQL database dump
--

-- Dumped from database version 15.13
-- Dumped by pg_dump version 15.13

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


--
-- Name: donation_status; Type: TYPE; Schema: public; Owner: microservice_user
--

CREATE TYPE public.donation_status AS ENUM (
    'completed',
    'pending',
    'failed',
    'refunded'
);


ALTER TYPE public.donation_status OWNER TO microservice_user;

--
-- Name: payment_method; Type: TYPE; Schema: public; Owner: microservice_user
--

CREATE TYPE public.payment_method AS ENUM (
    'MercadoPago',
    'Transferencia',
    'Efectivo',
    'Tarjeta'
);


ALTER TYPE public.payment_method OWNER TO microservice_user;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: activities; Type: TABLE; Schema: public; Owner: microservice_user
--

CREATE TABLE public.activities (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    campaign_id uuid NOT NULL,
    title character varying(255) NOT NULL,
    description text,
    date timestamp without time zone NOT NULL,
    type character varying(100) NOT NULL,
    author character varying(255) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.activities OWNER TO microservice_user;

--
-- Name: campaign_categories; Type: TABLE; Schema: public; Owner: microservice_user
--

CREATE TABLE public.campaign_categories (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying(100) NOT NULL,
    description text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.campaign_categories OWNER TO microservice_user;

--
-- Name: campaign_payment_methods; Type: TABLE; Schema: public; Owner: microservice_user
--

CREATE TABLE public.campaign_payment_methods (
    id integer NOT NULL,
    campaign_id uuid NOT NULL,
    payment_method_id integer NOT NULL,
    instructions text,
    is_active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.campaign_payment_methods OWNER TO microservice_user;

--
-- Name: campaign_payment_methods_id_seq; Type: SEQUENCE; Schema: public; Owner: microservice_user
--

CREATE SEQUENCE public.campaign_payment_methods_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.campaign_payment_methods_id_seq OWNER TO microservice_user;

--
-- Name: campaign_payment_methods_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: microservice_user
--

ALTER SEQUENCE public.campaign_payment_methods_id_seq OWNED BY public.campaign_payment_methods.id;


--
-- Name: campaigns; Type: TABLE; Schema: public; Owner: microservice_user
--

CREATE TABLE public.campaigns (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    title character varying(255) NOT NULL,
    description text NOT NULL,
    image character varying(500),
    goal numeric(12,2) NOT NULL,
    start_date timestamp with time zone NOT NULL,
    end_date timestamp with time zone NOT NULL,
    location character varying(255),
    category_id uuid,
    urgency integer,
    organizer_id uuid,
    status character varying(50) DEFAULT 'active'::character varying,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    beneficiary_name character varying(255),
    beneficiary_age integer,
    current_situation text,
    urgency_reason character varying(500),
    CONSTRAINT campaigns_beneficiary_age_check CHECK (((beneficiary_age IS NULL) OR (beneficiary_age >= 0))),
    CONSTRAINT campaigns_goal_check CHECK ((goal > (0)::numeric)),
    CONSTRAINT campaigns_urgency_check CHECK (((urgency >= 1) AND (urgency <= 10))),
    CONSTRAINT check_end_date_after_start CHECK ((end_date > start_date))
);


ALTER TABLE public.campaigns OWNER TO microservice_user;

--
-- Name: cash_locations; Type: TABLE; Schema: public; Owner: microservice_user
--

CREATE TABLE public.cash_locations (
    id integer NOT NULL,
    campaign_payment_method_id integer NOT NULL,
    location_name character varying(100) NOT NULL,
    address character varying(200) NOT NULL,
    contact_info character varying(100),
    available_hours character varying(100),
    additional_notes text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.cash_locations OWNER TO microservice_user;

--
-- Name: cash_locations_id_seq; Type: SEQUENCE; Schema: public; Owner: microservice_user
--

CREATE SEQUENCE public.cash_locations_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.cash_locations_id_seq OWNER TO microservice_user;

--
-- Name: cash_locations_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: microservice_user
--

ALTER SEQUENCE public.cash_locations_id_seq OWNED BY public.cash_locations.id;


--
-- Name: donations; Type: TABLE; Schema: public; Owner: microservice_user
--

CREATE TABLE public.donations (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    campaign_id uuid NOT NULL,
    donor_id uuid NOT NULL,
    amount numeric(10,2) NOT NULL,
    date timestamp with time zone NOT NULL,
    message text,
    is_anonymous boolean DEFAULT false,
    status public.donation_status DEFAULT 'pending'::public.donation_status NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    payment_method_id integer NOT NULL,
    CONSTRAINT donations_amount_check CHECK ((amount > (0)::numeric))
);


ALTER TABLE public.donations OWNER TO microservice_user;

--
-- Name: donors; Type: TABLE; Schema: public; Owner: microservice_user
--

CREATE TABLE public.donors (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    first_name character varying(255) NOT NULL,
    last_name character varying(255) NOT NULL,
    is_verified boolean DEFAULT false,
    phone character varying(50),
    email character varying(255) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.donors OWNER TO microservice_user;

--
-- Name: goose_db_version; Type: TABLE; Schema: public; Owner: microservice_user
--

CREATE TABLE public.goose_db_version (
    id integer NOT NULL,
    version_id bigint NOT NULL,
    is_applied boolean NOT NULL,
    tstamp timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.goose_db_version OWNER TO microservice_user;

--
-- Name: goose_db_version_id_seq; Type: SEQUENCE; Schema: public; Owner: microservice_user
--

ALTER TABLE public.goose_db_version ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.goose_db_version_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: organizers; Type: TABLE; Schema: public; Owner: microservice_user
--

CREATE TABLE public.organizers (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying(255) NOT NULL,
    avatar character varying(500),
    verified boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    user_id uuid,
    email character varying(255),
    phone character varying(20),
    website character varying(255),
    address character varying(500)
);


ALTER TABLE public.organizers OWNER TO microservice_user;

--
-- Name: payment_methods; Type: TABLE; Schema: public; Owner: microservice_user
--

CREATE TABLE public.payment_methods (
    id integer NOT NULL,
    code character varying(30) NOT NULL,
    name character varying(50) NOT NULL,
    is_active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.payment_methods OWNER TO microservice_user;

--
-- Name: payment_methods_id_seq; Type: SEQUENCE; Schema: public; Owner: microservice_user
--

CREATE SEQUENCE public.payment_methods_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.payment_methods_id_seq OWNER TO microservice_user;

--
-- Name: payment_methods_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: microservice_user
--

ALTER SEQUENCE public.payment_methods_id_seq OWNED BY public.payment_methods.id;


--
-- Name: permissions; Type: TABLE; Schema: public; Owner: microservice_user
--

CREATE TABLE public.permissions (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying(100) NOT NULL,
    resource character varying(50) NOT NULL,
    action character varying(20) NOT NULL,
    description text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.permissions OWNER TO microservice_user;

--
-- Name: receipts; Type: TABLE; Schema: public; Owner: microservice_user
--

CREATE TABLE public.receipts (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    campaign_id uuid NOT NULL,
    provider character varying(255) NOT NULL,
    name character varying(255) NOT NULL,
    description text,
    total numeric(12,2) NOT NULL,
    quantity integer DEFAULT 1,
    date timestamp with time zone NOT NULL,
    document_url character varying(500),
    note text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT receipts_quantity_check CHECK ((quantity >= 1)),
    CONSTRAINT receipts_total_check CHECK ((total > (0)::numeric))
);


ALTER TABLE public.receipts OWNER TO microservice_user;

--
-- Name: TABLE receipts; Type: COMMENT; Schema: public; Owner: microservice_user
--

COMMENT ON TABLE public.receipts IS 'Stores receipt documents for campaign expenses and transactions';


--
-- Name: COLUMN receipts.provider; Type: COMMENT; Schema: public; Owner: microservice_user
--

COMMENT ON COLUMN public.receipts.provider IS 'Name of the provider or vendor';


--
-- Name: COLUMN receipts.total; Type: COMMENT; Schema: public; Owner: microservice_user
--

COMMENT ON COLUMN public.receipts.total IS 'Total amount of the receipt';


--
-- Name: COLUMN receipts.quantity; Type: COMMENT; Schema: public; Owner: microservice_user
--

COMMENT ON COLUMN public.receipts.quantity IS 'Number of items or units';


--
-- Name: COLUMN receipts.document_url; Type: COMMENT; Schema: public; Owner: microservice_user
--

COMMENT ON COLUMN public.receipts.document_url IS 'URL to the PDF document stored in S3';


--
-- Name: COLUMN receipts.note; Type: COMMENT; Schema: public; Owner: microservice_user
--

COMMENT ON COLUMN public.receipts.note IS 'Additional notes or comments about the receipt';


--
-- Name: role_permissions; Type: TABLE; Schema: public; Owner: microservice_user
--

CREATE TABLE public.role_permissions (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    role_id uuid NOT NULL,
    permission_id uuid NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.role_permissions OWNER TO microservice_user;

--
-- Name: roles; Type: TABLE; Schema: public; Owner: microservice_user
--

CREATE TABLE public.roles (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying(50) NOT NULL,
    description text,
    is_active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.roles OWNER TO microservice_user;

--
-- Name: transfer_details; Type: TABLE; Schema: public; Owner: microservice_user
--

CREATE TABLE public.transfer_details (
    id integer NOT NULL,
    campaign_payment_method_id integer NOT NULL,
    bank_name character varying(100) NOT NULL,
    account_holder character varying(100) NOT NULL,
    cbu character varying(22) NOT NULL,
    alias character varying(30),
    swift_code character varying(11),
    additional_notes text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.transfer_details OWNER TO microservice_user;

--
-- Name: transfer_details_id_seq; Type: SEQUENCE; Schema: public; Owner: microservice_user
--

CREATE SEQUENCE public.transfer_details_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.transfer_details_id_seq OWNER TO microservice_user;

--
-- Name: transfer_details_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: microservice_user
--

ALTER SEQUENCE public.transfer_details_id_seq OWNED BY public.transfer_details.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: microservice_user
--

CREATE TABLE public.users (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    email character varying(255) NOT NULL,
    password_hash character varying(255) NOT NULL,
    first_name character varying(100),
    last_name character varying(100),
    is_active boolean DEFAULT true,
    is_verified boolean DEFAULT false,
    reset_token character varying(255),
    reset_token_expires_at timestamp with time zone,
    last_login timestamp with time zone,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    role_id uuid NOT NULL
);


ALTER TABLE public.users OWNER TO microservice_user;

--
-- Name: campaign_payment_methods id; Type: DEFAULT; Schema: public; Owner: microservice_user
--

ALTER TABLE ONLY public.campaign_payment_methods ALTER COLUMN id SET DEFAULT nextval('public.campaign_payment_methods_id_seq'::regclass);


--
-- Name: cash_locations id; Type: DEFAULT; Schema: public; Owner: microservice_user
--

ALTER TABLE ONLY public.cash_locations ALTER COLUMN id SET DEFAULT nextval('public.cash_locations_id_seq'::regclass);


--
-- Name: payment_methods id; Type: DEFAULT; Schema: public; Owner: microservice_user
--

ALTER TABLE ONLY public.payment_methods ALTER COLUMN id SET DEFAULT nextval('public.payment_methods_id_seq'::regclass);


--
-- Name: transfer_details id; Type: DEFAULT; Schema: public; Owner: microservice_user
--

ALTER TABLE ONLY public.transfer_details ALTER COLUMN id SET DEFAULT nextval('public.transfer_details_id_seq'::regclass);


--
-- Name: activities activities_pkey; Type: CONSTRAINT; Schema: public; Owner: microservice_user
--

ALTER TABLE ONLY public.activities
    ADD CONSTRAINT activities_pkey PRIMARY KEY (id);


--
-- Name: campaign_categories campaign_categories_name_key; Type: CONSTRAINT; Schema: public; Owner: microservice_user
--

ALTER TABLE ONLY public.campaign_categories
    ADD CONSTRAINT campaign_categories_name_key UNIQUE (name);


--
-- Name: campaign_categories campaign_categories_pkey; Type: CONSTRAINT; Schema: public; Owner: microservice_user
--

ALTER TABLE ONLY public.campaign_categories
    ADD CONSTRAINT campaign_categories_pkey PRIMARY KEY (id);


--
-- Name: campaign_payment_methods campaign_payment_methods_campaign_id_payment_method_id_key; Type: CONSTRAINT; Schema: public; Owner: microservice_user
--

ALTER TABLE ONLY public.campaign_payment_methods
    ADD CONSTRAINT campaign_payment_methods_campaign_id_payment_method_id_key UNIQUE (campaign_id, payment_method_id);


--
-- Name: campaign_payment_methods campaign_payment_methods_pkey; Type: CONSTRAINT; Schema: public; Owner: microservice_user
--

ALTER TABLE ONLY public.campaign_payment_methods
    ADD CONSTRAINT campaign_payment_methods_pkey PRIMARY KEY (id);


--
-- Name: campaigns campaigns_pkey; Type: CONSTRAINT; Schema: public; Owner: microservice_user
--

ALTER TABLE ONLY public.campaigns
    ADD CONSTRAINT campaigns_pkey PRIMARY KEY (id);


--
-- Name: cash_locations cash_locations_pkey; Type: CONSTRAINT; Schema: public; Owner: microservice_user
--

ALTER TABLE ONLY public.cash_locations
    ADD CONSTRAINT cash_locations_pkey PRIMARY KEY (id);


--
-- Name: donations donations_pkey; Type: CONSTRAINT; Schema: public; Owner: microservice_user
--

ALTER TABLE ONLY public.donations
    ADD CONSTRAINT donations_pkey PRIMARY KEY (id);


--
-- Name: donors donors_email_key; Type: CONSTRAINT; Schema: public; Owner: microservice_user
--

ALTER TABLE ONLY public.donors
    ADD CONSTRAINT donors_email_key UNIQUE (email);


--
-- Name: donors donors_pkey; Type: CONSTRAINT; Schema: public; Owner: microservice_user
--

ALTER TABLE ONLY public.donors
    ADD CONSTRAINT donors_pkey PRIMARY KEY (id);


--
-- Name: goose_db_version goose_db_version_pkey; Type: CONSTRAINT; Schema: public; Owner: microservice_user
--

ALTER TABLE ONLY public.goose_db_version
    ADD CONSTRAINT goose_db_version_pkey PRIMARY KEY (id);


--
-- Name: organizers organizers_pkey; Type: CONSTRAINT; Schema: public; Owner: microservice_user
--

ALTER TABLE ONLY public.organizers
    ADD CONSTRAINT organizers_pkey PRIMARY KEY (id);


--
-- Name: payment_methods payment_methods_code_key; Type: CONSTRAINT; Schema: public; Owner: microservice_user
--

ALTER TABLE ONLY public.payment_methods
    ADD CONSTRAINT payment_methods_code_key UNIQUE (code);


--
-- Name: payment_methods payment_methods_pkey; Type: CONSTRAINT; Schema: public; Owner: microservice_user
--

ALTER TABLE ONLY public.payment_methods
    ADD CONSTRAINT payment_methods_pkey PRIMARY KEY (id);


--
-- Name: permissions permissions_name_key; Type: CONSTRAINT; Schema: public; Owner: microservice_user
--

ALTER TABLE ONLY public.permissions
    ADD CONSTRAINT permissions_name_key UNIQUE (name);


--
-- Name: permissions permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: microservice_user
--

ALTER TABLE ONLY public.permissions
    ADD CONSTRAINT permissions_pkey PRIMARY KEY (id);


--
-- Name: receipts receipts_pkey; Type: CONSTRAINT; Schema: public; Owner: microservice_user
--

ALTER TABLE ONLY public.receipts
    ADD CONSTRAINT receipts_pkey PRIMARY KEY (id);


--
-- Name: role_permissions role_permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: microservice_user
--

ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT role_permissions_pkey PRIMARY KEY (id);


--
-- Name: role_permissions role_permissions_role_id_permission_id_key; Type: CONSTRAINT; Schema: public; Owner: microservice_user
--

ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT role_permissions_role_id_permission_id_key UNIQUE (role_id, permission_id);


--
-- Name: roles roles_name_key; Type: CONSTRAINT; Schema: public; Owner: microservice_user
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_name_key UNIQUE (name);


--
-- Name: roles roles_pkey; Type: CONSTRAINT; Schema: public; Owner: microservice_user
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_pkey PRIMARY KEY (id);


--
-- Name: transfer_details transfer_details_pkey; Type: CONSTRAINT; Schema: public; Owner: microservice_user
--

ALTER TABLE ONLY public.transfer_details
    ADD CONSTRAINT transfer_details_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: microservice_user
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: microservice_user
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_activities_campaign_id; Type: INDEX; Schema: public; Owner: microservice_user
--

CREATE INDEX idx_activities_campaign_id ON public.activities USING btree (campaign_id);


--
-- Name: idx_activities_date; Type: INDEX; Schema: public; Owner: microservice_user
--

CREATE INDEX idx_activities_date ON public.activities USING btree (date);


--
-- Name: idx_campaign_payment_methods_campaign_id; Type: INDEX; Schema: public; Owner: microservice_user
--

CREATE INDEX idx_campaign_payment_methods_campaign_id ON public.campaign_payment_methods USING btree (campaign_id);


--
-- Name: idx_campaign_payment_methods_payment_method_id; Type: INDEX; Schema: public; Owner: microservice_user
--

CREATE INDEX idx_campaign_payment_methods_payment_method_id ON public.campaign_payment_methods USING btree (payment_method_id);


--
-- Name: idx_campaigns_beneficiary_age; Type: INDEX; Schema: public; Owner: microservice_user
--

CREATE INDEX idx_campaigns_beneficiary_age ON public.campaigns USING btree (beneficiary_age);


--
-- Name: idx_campaigns_beneficiary_name; Type: INDEX; Schema: public; Owner: microservice_user
--

CREATE INDEX idx_campaigns_beneficiary_name ON public.campaigns USING btree (beneficiary_name);


--
-- Name: idx_campaigns_category_id; Type: INDEX; Schema: public; Owner: microservice_user
--

CREATE INDEX idx_campaigns_category_id ON public.campaigns USING btree (category_id);


--
-- Name: idx_campaigns_end_date; Type: INDEX; Schema: public; Owner: microservice_user
--

CREATE INDEX idx_campaigns_end_date ON public.campaigns USING btree (end_date);


--
-- Name: idx_campaigns_organizer_id; Type: INDEX; Schema: public; Owner: microservice_user
--

CREATE INDEX idx_campaigns_organizer_id ON public.campaigns USING btree (organizer_id);


--
-- Name: idx_campaigns_start_date; Type: INDEX; Schema: public; Owner: microservice_user
--

CREATE INDEX idx_campaigns_start_date ON public.campaigns USING btree (start_date);


--
-- Name: idx_campaigns_status; Type: INDEX; Schema: public; Owner: microservice_user
--

CREATE INDEX idx_campaigns_status ON public.campaigns USING btree (status);


--
-- Name: idx_cash_locations_campaign_payment_method_id; Type: INDEX; Schema: public; Owner: microservice_user
--

CREATE INDEX idx_cash_locations_campaign_payment_method_id ON public.cash_locations USING btree (campaign_payment_method_id);


--
-- Name: idx_donations_campaign_id; Type: INDEX; Schema: public; Owner: microservice_user
--

CREATE INDEX idx_donations_campaign_id ON public.donations USING btree (campaign_id);


--
-- Name: idx_donations_date; Type: INDEX; Schema: public; Owner: microservice_user
--

CREATE INDEX idx_donations_date ON public.donations USING btree (date);


--
-- Name: idx_donations_donor_id; Type: INDEX; Schema: public; Owner: microservice_user
--

CREATE INDEX idx_donations_donor_id ON public.donations USING btree (donor_id);


--
-- Name: idx_donations_status; Type: INDEX; Schema: public; Owner: microservice_user
--

CREATE INDEX idx_donations_status ON public.donations USING btree (status);


--
-- Name: idx_donors_email; Type: INDEX; Schema: public; Owner: microservice_user
--

CREATE INDEX idx_donors_email ON public.donors USING btree (email);


--
-- Name: idx_organizers_user_id; Type: INDEX; Schema: public; Owner: microservice_user
--

CREATE INDEX idx_organizers_user_id ON public.organizers USING btree (user_id);


--
-- Name: idx_permissions_resource_action; Type: INDEX; Schema: public; Owner: microservice_user
--

CREATE INDEX idx_permissions_resource_action ON public.permissions USING btree (resource, action);


--
-- Name: idx_receipts_campaign_id; Type: INDEX; Schema: public; Owner: microservice_user
--

CREATE INDEX idx_receipts_campaign_id ON public.receipts USING btree (campaign_id);


--
-- Name: idx_receipts_date; Type: INDEX; Schema: public; Owner: microservice_user
--

CREATE INDEX idx_receipts_date ON public.receipts USING btree (date);


--
-- Name: idx_receipts_provider; Type: INDEX; Schema: public; Owner: microservice_user
--

CREATE INDEX idx_receipts_provider ON public.receipts USING btree (provider);


--
-- Name: idx_role_permissions_permission_id; Type: INDEX; Schema: public; Owner: microservice_user
--

CREATE INDEX idx_role_permissions_permission_id ON public.role_permissions USING btree (permission_id);


--
-- Name: idx_role_permissions_role_id; Type: INDEX; Schema: public; Owner: microservice_user
--

CREATE INDEX idx_role_permissions_role_id ON public.role_permissions USING btree (role_id);


--
-- Name: idx_roles_active; Type: INDEX; Schema: public; Owner: microservice_user
--

CREATE INDEX idx_roles_active ON public.roles USING btree (is_active);


--
-- Name: idx_roles_name; Type: INDEX; Schema: public; Owner: microservice_user
--

CREATE INDEX idx_roles_name ON public.roles USING btree (name);


--
-- Name: idx_transfer_details_campaign_payment_method_id; Type: INDEX; Schema: public; Owner: microservice_user
--

CREATE INDEX idx_transfer_details_campaign_payment_method_id ON public.transfer_details USING btree (campaign_payment_method_id);


--
-- Name: idx_users_email; Type: INDEX; Schema: public; Owner: microservice_user
--

CREATE INDEX idx_users_email ON public.users USING btree (email);


--
-- Name: idx_users_reset_token; Type: INDEX; Schema: public; Owner: microservice_user
--

CREATE INDEX idx_users_reset_token ON public.users USING btree (reset_token);


--
-- Name: idx_users_role_id; Type: INDEX; Schema: public; Owner: microservice_user
--

CREATE INDEX idx_users_role_id ON public.users USING btree (role_id);


--
-- Name: activities activities_campaign_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: microservice_user
--

ALTER TABLE ONLY public.activities
    ADD CONSTRAINT activities_campaign_id_fkey FOREIGN KEY (campaign_id) REFERENCES public.campaigns(id) ON DELETE CASCADE;


--
-- Name: campaign_payment_methods campaign_payment_methods_campaign_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: microservice_user
--

ALTER TABLE ONLY public.campaign_payment_methods
    ADD CONSTRAINT campaign_payment_methods_campaign_id_fkey FOREIGN KEY (campaign_id) REFERENCES public.campaigns(id) ON DELETE CASCADE;


--
-- Name: campaign_payment_methods campaign_payment_methods_payment_method_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: microservice_user
--

ALTER TABLE ONLY public.campaign_payment_methods
    ADD CONSTRAINT campaign_payment_methods_payment_method_id_fkey FOREIGN KEY (payment_method_id) REFERENCES public.payment_methods(id) ON DELETE CASCADE;


--
-- Name: campaigns campaigns_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: microservice_user
--

ALTER TABLE ONLY public.campaigns
    ADD CONSTRAINT campaigns_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.campaign_categories(id);


--
-- Name: campaigns campaigns_organizer_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: microservice_user
--

ALTER TABLE ONLY public.campaigns
    ADD CONSTRAINT campaigns_organizer_id_fkey FOREIGN KEY (organizer_id) REFERENCES public.organizers(id);


--
-- Name: cash_locations cash_locations_campaign_payment_method_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: microservice_user
--

ALTER TABLE ONLY public.cash_locations
    ADD CONSTRAINT cash_locations_campaign_payment_method_id_fkey FOREIGN KEY (campaign_payment_method_id) REFERENCES public.campaign_payment_methods(id) ON DELETE CASCADE;


--
-- Name: donations donations_campaign_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: microservice_user
--

ALTER TABLE ONLY public.donations
    ADD CONSTRAINT donations_campaign_id_fkey FOREIGN KEY (campaign_id) REFERENCES public.campaigns(id);


--
-- Name: donations donations_donor_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: microservice_user
--

ALTER TABLE ONLY public.donations
    ADD CONSTRAINT donations_donor_id_fkey FOREIGN KEY (donor_id) REFERENCES public.donors(id);


--
-- Name: donations fk_donations_payment_method; Type: FK CONSTRAINT; Schema: public; Owner: microservice_user
--

ALTER TABLE ONLY public.donations
    ADD CONSTRAINT fk_donations_payment_method FOREIGN KEY (payment_method_id) REFERENCES public.payment_methods(id);


--
-- Name: organizers fk_organizers_user_id; Type: FK CONSTRAINT; Schema: public; Owner: microservice_user
--

ALTER TABLE ONLY public.organizers
    ADD CONSTRAINT fk_organizers_user_id FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: receipts receipts_campaign_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: microservice_user
--

ALTER TABLE ONLY public.receipts
    ADD CONSTRAINT receipts_campaign_id_fkey FOREIGN KEY (campaign_id) REFERENCES public.campaigns(id) ON DELETE CASCADE;


--
-- Name: role_permissions role_permissions_permission_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: microservice_user
--

ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT role_permissions_permission_id_fkey FOREIGN KEY (permission_id) REFERENCES public.permissions(id) ON DELETE CASCADE;


--
-- Name: role_permissions role_permissions_role_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: microservice_user
--

ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT role_permissions_role_id_fkey FOREIGN KEY (role_id) REFERENCES public.roles(id) ON DELETE CASCADE;


--
-- Name: transfer_details transfer_details_campaign_payment_method_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: microservice_user
--

ALTER TABLE ONLY public.transfer_details
    ADD CONSTRAINT transfer_details_campaign_payment_method_id_fkey FOREIGN KEY (campaign_payment_method_id) REFERENCES public.campaign_payment_methods(id) ON DELETE CASCADE;


--
-- Name: users users_role_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: microservice_user
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_role_id_fkey FOREIGN KEY (role_id) REFERENCES public.roles(id);


--
-- PostgreSQL database dump complete
--

