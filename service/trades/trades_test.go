package trades_test

import (
	tradeRepositoryMock "VincentLimarus/stock-analyzer-performance/mocks/repository/trades"
	db "VincentLimarus/stock-analyzer-performance/model/db/trades"
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	svcTrade "VincentLimarus/stock-analyzer-performance/service/trades"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	mockRepository "VincentLimarus/stock-analyzer-performance/mocks/repository"
)

type TradeTestSuite struct {
	suite.Suite
	mocks mockTrade
}

type SqlxDB *sqlx.DB
type SqlMock sqlmock.Sqlmock

type mockTrade struct {
	db SqlxDB
	mockDB SqlMock
	mockRepository *mockRepository.IRegistry
	tradeRepository *tradeRepositoryMock.ITrade
	svc *svcTrade.Trade
}

func (suite *TradeTestSuite) SetupTest() {
	db, mockDb, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		suite.FailNow(err.Error())
	}

	// Mock Repo Registry 
	mockRepoRegistry := mockRepository.NewIRegistry(suite.T())

	// Mock Repository 
	mockTradeRepo := tradeRepositoryMock.NewITrade(suite.T())

	svc := svcTrade.NewTrade(mockRepoRegistry)

	suite.mocks = mockTrade{
		db : sqlx.NewDb(db, "sqlmock"),
		mockDB : mockDb,
		mockRepository: mockRepoRegistry,
		tradeRepository: mockTradeRepo,
		svc: svc,
	}
}

func (suite *TradeTestSuite) TearDownTest() {
	fmt.Println("teardown")
}

func (suite *TradeTestSuite) TestGetListTradeByYears() {
	type (
		args struct {
			ctx context.Context
			years []string
		}
		mockSetupFunc func (m mockTrade, args args)
		wantErr bool
		expectedError error
		testCase      struct {
			name          string
			args          args
			mockSetupFunc mockSetupFunc
			wantErr       wantErr
			expectedError expectedError
		}
	)
	
	listTradeData := []db.Trade{
		{
			SerialCode: "SC001",
			Status: "CLOSED",
			StockCode: "STK001",
			Lot: 100,
			AvgBuyPrice: 10.5,
			AvgSellPrice: 12.0,
			DividenAmount: nil,
		},
	}

	tests := []testCase{	
		{
			name: "should return error when repository returns error",
			args: args{
				ctx: context.Background(),
				years: []string{"2022", "2023"},
			},
			mockSetupFunc: func(m mockTrade, args args) {
				m.mockRepository.On("GetTrade").Return(m.tradeRepository).Once()
				m.tradeRepository.On("GetListTradeByYears", args.ctx, args.years).
					Return(nil, fmt.Errorf("repository error")).Once()
			},
			wantErr: true,
			expectedError: fmt.Errorf("repository error"),	
		},
		{
			name: "should return no error when repository returns data successfully",
			args: args{
				ctx: context.Background(),
				years: []string{"2022", "2023"},
			},
			mockSetupFunc: func(m mockTrade, args args) {
				m.mockRepository.On("GetTrade").Return(m.tradeRepository).Once()
				m.tradeRepository.On("GetListTradeByYears", args.ctx, args.years).
					Return(listTradeData, nil).Once()
			},
			wantErr: false,
			expectedError: nil,	
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			if tt.mockSetupFunc != nil {
				tt.mockSetupFunc(suite.mocks, tt.args)
			}

			_, err := suite.mocks.svc.GetListTradeByYears(tt.args.ctx, tt.args.years)
			if tt.wantErr {
				assert.Error(suite.T(), err)
				if tt.expectedError != nil {
					assert.Contains(suite.T(), err.Error(), tt.expectedError.Error())
				}
			} else {
				assert.NoError(suite.T(), err)
			}
		})
	}
} 	

func TestTradeSuite(t *testing.T) {
	suite.Run(t, new(TradeTestSuite))
}
