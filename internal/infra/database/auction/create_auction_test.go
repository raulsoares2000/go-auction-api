package auction

import (
	"context"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// setupTestDB é uma função auxiliar para configurar um banco de dados de teste.
func setupTestDB(t *testing.T) (*mongo.Database, func()) {
	t.Helper()

	// Usamos uma string de conexão para um MongoDB local, ideal para testes
	// Pode ser o mesmo MongoDB que você está a usar com o Docker
	mongoURL := "mongodb://admin:admin@localhost:27017"
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURL))
	if err != nil {
		t.Fatalf("Falha ao conectar ao MongoDB: %v", err)
	}

	db := client.Database("auctions_test") // Usamos uma base de dados separada para testes

	// teardown function: será chamada no final do teste para limpar tudo
	teardown := func() {
		db.Drop(context.Background()) // Apaga a base de dados de teste
		client.Disconnect(context.Background())
	}

	return db, teardown
}

func TestAuctionRepository_CreateAuction_ShouldCloseAuction_AfterInterval(t *testing.T) {
	// --- 1. Configuração (Arrange) ---
	db, teardown := setupTestDB(t)
	defer teardown()

	repo := NewAuctionRepository(db)

	auctionInterval := "5s"
	os.Setenv("AUCTION_INTERVAL", auctionInterval)

	auction, _ := auction_entity.CreateAuction(
		"Test Product",
		"Test Category",
		"Test description for product",
		auction_entity.New,
	)

	// --- 2. Execução (Act) ---
	err := repo.CreateAuction(context.Background(), auction)
	assert.Nil(t, err)

	// --- 3. Verificação Imediata (Assert Initial State) ---
	// Buscamos o leilão no banco de dados IMEDIATAMENTE após a criação.
	var initialAuctionMongo AuctionEntityMongo
	filter := bson.M{"_id": auction.Id}
	errDbInitial := repo.Collection.FindOne(context.Background(), filter).Decode(&initialAuctionMongo)

	// Verificamos se ele foi encontrado e se o status inicial está correto.
	assert.Nil(t, errDbInitial, "Erro ao buscar o leilão imediatamente após a criação")
	assert.Equal(t, auction_entity.Active, initialAuctionMongo.Status, "O leilão deveria ser criado com o status 'Active'")

	// --- 4. Aguardar a expiração do leilão ---
	// Fazemos a pausa para dar tempo à goroutine de executar a atualização.
	waitTime, _ := time.ParseDuration(auctionInterval)
	time.Sleep(waitTime + 1*time.Second)

	// --- 5. Verificação Final (Assert Final State) ---
	// Buscamos o MESMO leilão novamente, agora para verificar o status final.
	var finalAuctionMongo AuctionEntityMongo
	errDbFinal := repo.Collection.FindOne(context.Background(), filter).Decode(&finalAuctionMongo)

	// Verificamos se ele foi encontrado e se o status foi atualizado para 'Completed'.
	assert.Nil(t, errDbFinal, "Erro ao buscar o leilão após o tempo de expiração")
	assert.Equal(t, auction_entity.Completed, finalAuctionMongo.Status, "O status do leilão deveria ser 'Completed' após o intervalo")
}
